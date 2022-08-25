package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

var (
	numReqs  uint64
	lastReqs uint64
)

func newServer(opts serverOptions) (server, error) {
	s := server{
		opts:      opts,
		userStore: userStore{pg: opts.pg},
		urlStore:  urlStore{pg: opts.pg, rds: opts.rds},
		idStore:   idStore{rds: opts.rds},
	}

	return s, nil
}

type server struct {
	idStore
	urlStore
	userStore

	opts serverOptions
}

func (s server) startCountingReqs() {
	for {
		time.Sleep(2 * time.Second)
		fmt.Printf("%d reqs/s\n", (numReqs-lastReqs)/2)

		lastReqs = numReqs
	}
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		switch r.URL.Path {
		case "/oauth/callback":
			s.handleCallback(w, r)
			return
		case "/oauth/login":
			s.handleLogin(w, r)
			return
		case "/shorturls":
			s.handleListShortUrls(w, r)

		default:
			s.handleRedirectShortLink(w, r)
			return
		}
	}

	if r.Method == "POST" {
		s.handleCreateShortLink(w, r)
	}
}

func (s server) handleListShortUrls(w http.ResponseWriter, r *http.Request) {
	u, err := s.getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	urls, err := s.urlStore.listUrlsByUserId(u.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(urls)
}

func (s server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if s.opts.oauthClientID == "" {
		http.Error(w, "no oauth client id", http.StatusInternalServerError)
		return
	}

	v := url.Values{}
	v.Add("client_id", s.opts.oauthClientID)
	v.Add("redirect_uri", s.opts.oauthRedirectURL)

	authURL, _ := url.Parse(s.opts.oauthAuthURL)
	authURL.RawQuery = v.Encode()

	http.Redirect(w, r, authURL.String(), http.StatusMovedPermanently)
}

func (s server) handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query()["code"]

	tokenParams := url.Values{}
	tokenParams.Add("client_id", s.opts.oauthClientID)
	tokenParams.Add("client_secret", s.opts.oauthClientSecret)
	tokenParams.Add("code", code[0])

	tokenResp, err := http.PostForm(s.opts.oauthTokenURL, tokenParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer tokenResp.Body.Close()

	tokenRespBytes, err := ioutil.ReadAll(tokenResp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Reusing variable
	tokenParams, err = url.ParseQuery(string(tokenRespBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, ok := tokenParams["access_token"]
	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userReq, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	userReq.Header.Add("Authorization", "token "+accessToken[0])

	userResp, err := http.DefaultClient.Do(userReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer userResp.Body.Close()

	u := user{
		AccessToken: accessToken[0],
	}

	json.NewDecoder(userResp.Body).Decode(&u)

	err = s.userStore.setUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(tokenRespBytes)
}

func (s server) handleRedirectShortLink(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(&numReqs, 1)

	id := r.RequestURI[1:]
	if len(id) != 4 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	url, err := s.urlStore.getUrl(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func (s server) handleCreateShortLink(w http.ResponseWriter, r *http.Request) {
	u, err := s.getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	req := struct {
		URL string `json:"url"`
	}{}

	res := struct {
		URL string `json:"url"`
	}{}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "could not decode request", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	id, err := s.idStore.getNextID()
	if err != nil {
		http.Error(w, "could not get next id", http.StatusInternalServerError)
		return
	}

	err = s.urlStore.setUrl(u.ID, id, req.URL)
	if err != nil {
		http.Error(w, "could not get store url", http.StatusInternalServerError)
		return
	}

	res.URL = s.opts.protocol + "://" + s.opts.hostname + "/" + id
	json.NewEncoder(w).Encode(res)
}

func (s server) getUser(r *http.Request) (user, error) {
	authorization, ok := r.Header["Authorization"]
	if !ok {
		return user{}, fmt.Errorf("missing authorization header")
	}

	accessToken := authorization[0][6:]

	return s.userStore.getUserByAccessToken(accessToken)
}
