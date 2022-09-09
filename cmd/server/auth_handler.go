package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type authHandler struct {
	opts      options
	userStore userStore
}

func newAuthHandler(opts options) authHandler {
	return authHandler{
		userStore: userStore{
			opts.pg,
		},

		opts: opts,
	}
}

func (s authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqeustsCounter.Inc()

	if r.Method == "GET" {
		switch r.URL.Path {
		case "/auth/callback":
			s.handleCallback(w, r)
			return
		case "/auth/login":
			s.handleLogin(w, r)
			return
		}
	}

	http.Error(w, "", http.StatusBadRequest)
}

func (s authHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
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

func (s authHandler) handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query()["code"]

	tokenParams := url.Values{}
	tokenParams.Add("client_id", s.opts.oauthClientID)
	tokenParams.Add("client_secret", s.opts.oauthClientSecret)
	tokenParams.Add("code", code[0])

	tokenResp, err := http.PostForm(s.opts.oauthTokenURL, tokenParams)
	if err != nil {
		errorCounter.Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer tokenResp.Body.Close()

	tokenRespBytes, err := ioutil.ReadAll(tokenResp.Body)
	if err != nil {
		errorCounter.Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Reusing variable
	tokenParams, err = url.ParseQuery(string(tokenRespBytes))
	if err != nil {
		errorCounter.Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, ok := tokenParams["access_token"]
	if !ok {
		errorCounter.Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userReq, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	userReq.Header.Add("Authorization", "token "+accessToken[0])

	userResp, err := http.DefaultClient.Do(userReq)
	if err != nil {
		errorCounter.Inc()
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
		errorCounter.Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(tokenRespBytes)
}
