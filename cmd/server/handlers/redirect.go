package handlers

import (
	"net/http"

	"github.com/nwehr/shorty/cmd/server/models"
	"github.com/nwehr/shorty/cmd/server/options"
)

func Redirect(opts options.Options) http.HandlerFunc {
	urlStore := models.URLMapStore{
		PG:  opts.PG,
		RDS: opts.RDS,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		key := r.RequestURI[1:]

		m, err := urlStore.Get(key)
		if err != nil {
			opts.Logger.Debugf("could not find url by key %s\n", key)

			errorCounter.Inc()
			w.WriteHeader(http.StatusNotFound)
			return
		}

		opts.Logger.Debugf("redirecting %s -> %s\n", key, m.LongURL)

		http.Redirect(w, r, m.LongURL, http.StatusTemporaryRedirect)
	}
}
