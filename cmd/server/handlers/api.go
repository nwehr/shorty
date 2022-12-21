package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nwehr/shorty/cmd/server/models"
	"github.com/nwehr/shorty/cmd/server/options"
)

type ApiHandler struct {
	urlStore models.URLMapStore
	opts     options.Options
}

func NewAPIHandler(opts options.Options) ApiHandler {
	return ApiHandler{
		urlStore: models.URLMapStore{
			PG:  opts.PG,
			RDS: opts.RDS,
		},

		opts: opts,
	}
}

func (s ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqeustsCounter.Inc()

	if r.Method == "GET" && r.URL.Path == "/api/list" {
		s.handleListShortUrls(w, r)
		return
	}

	if r.Method == "POST" && r.URL.Path == "/api/create" {
		s.handleCreateShortLink(w, r)
		return
	}

	s.opts.Logger.Errorf("%s %s is a bad request\n", r.Method, r.URL.Path)

	http.Error(w, "", http.StatusBadRequest)
}

func (s ApiHandler) handleListShortUrls(w http.ResponseWriter, r *http.Request) {
	authorization, err := getAuthorization(r)
	if err != nil {
		errorCounter.Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	urls, err := s.urlStore.List(AccessToken(authorization.AccessToken).GetIssuer())
	if err != nil {
		s.opts.Logger.Errorf("could not get url: %v\n", err)

		errorCounter.Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(urls)
}

func (s ApiHandler) handleCreateShortLink(w http.ResponseWriter, r *http.Request) {
	authorization, err := getAuthorization(r)
	if err != nil {
		errorCounter.Inc()
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	m := models.URLMap{}

	err = json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		s.opts.Logger.Errorf("could not decode request: %v\n", err)

		errorCounter.Inc()
		http.Error(w, "could not decode request", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	m, err = s.urlStore.Set(AccessToken(authorization.AccessToken).GetIssuer(), m)
	if err != nil {
		s.opts.Logger.Errorf("could not set url: %v\n", err)

		errorCounter.Inc()
		http.Error(w, "could not set store url", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(m)
}
