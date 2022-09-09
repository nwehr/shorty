package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type apiHandler struct {
	urlStore  urlStore
	keyStore  keyStore
	userStore userStore
	opts      options
}

func newAPIHandler(opts options) apiHandler {
	return apiHandler{
		urlStore: urlStore{
			pg:  opts.pg,
			rds: opts.rds,
		},

		keyStore: keyStore{
			rds: opts.rds,
		},

		userStore: userStore{
			pg: opts.pg,
		},

		opts: opts,
	}
}

func (s apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqeustsCounter.Inc()

	if r.Method == "GET" && r.URL.Path == "/api/list" {
		s.handleListShortUrls(w, r)
		return
	}

	if r.Method == "POST" && r.URL.Path == "/api/create" {
		s.handleCreateShortLink(w, r)
		return
	}

	s.opts.logger.Errorf("%s %s is a bad request\n", r.Method, r.URL.Path)

	http.Error(w, "", http.StatusBadRequest)
}

func (s apiHandler) handleListShortUrls(w http.ResponseWriter, r *http.Request) {
	u, err := s.getUser(r)
	if err != nil {
		s.opts.logger.Errorf("could not get user: %v\n", err)

		errorCounter.Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	urls, err := s.urlStore.listURLsByUserID(u.ID, s.opts.publicURL)
	if err != nil {
		s.opts.logger.Errorf("could not get url: %v\n", err)

		errorCounter.Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(urls)
}

func (s apiHandler) handleCreateShortLink(w http.ResponseWriter, r *http.Request) {
	u, err := s.getUser(r)
	if err != nil {
		s.opts.logger.Errorf("could not get user: %v\n", err)

		errorCounter.Inc()
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	mapped := mappedURL{}

	err = json.NewDecoder(r.Body).Decode(&mapped)
	if err != nil {
		s.opts.logger.Errorf("could not decode request: %v\n", err)

		errorCounter.Inc()
		http.Error(w, "could not decode request", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	key, err := s.keyStore.getNextKey()
	if err != nil {
		s.opts.logger.Errorf("could not get next key: %v\n", err)

		errorCounter.Inc()
		http.Error(w, "could not get next id", http.StatusInternalServerError)
		return
	}

	err = s.urlStore.setURL(u.ID, key, mapped.LongURL)
	if err != nil {
		s.opts.logger.Errorf("could not set url: %v\n", err)

		errorCounter.Inc()
		http.Error(w, "could not set store url", http.StatusInternalServerError)
		return
	}

	mapped.ShortURL = s.opts.publicURL + "/" + key

	json.NewEncoder(w).Encode(mapped)
}

func (s apiHandler) getUser(r *http.Request) (user, error) {
	authorization, ok := r.Header["Authorization"]
	if !ok {
		return user{}, fmt.Errorf("missing authorization header")
	}

	accessToken := authorization[0][6:]

	return s.userStore.getUserByAccessToken(accessToken)
}
