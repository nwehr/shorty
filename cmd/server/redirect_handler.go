package main

import "net/http"

func redirectHandler(opts options) http.HandlerFunc {
	urlStore := urlStore{
		pg:  opts.pg,
		rds: opts.rds,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		key := r.RequestURI[1:]

		url, err := urlStore.getURL(key)
		if err != nil {
			opts.logger.Debugf("could not find url by key %s\n", key)

			errorCounter.Inc()
			w.WriteHeader(http.StatusNotFound)
			return
		}

		opts.logger.Debugf("redirecting %s -> %s\n", key, url)

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
