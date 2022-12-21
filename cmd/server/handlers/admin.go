package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/nwehr/shorty/cmd/server/models"
	"github.com/nwehr/shorty/cmd/server/options"
)

type (
	IndexPage struct {
		Opts     options.Options
		URLs     []models.URLMap
		LoggedIn bool
		Username string
		Issuer   string
	}
)

func Admin(opts options.Options) http.HandlerFunc {
	store := models.URLMapStore{
		PG:  opts.PG,
		RDS: opts.RDS,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("templates/index.gohtml")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		authorization, err := getAuthorization(r)
		if err != nil {
			fmt.Println(err)
		}

		p := IndexPage{
			Opts:     opts,
			LoggedIn: err == nil,
			Issuer:   AccessToken(authorization.AccessToken).GetIssuer(),
			Username: AccessToken(authorization.AccessToken).GetUsername(),
		}

		if err == nil {
			p.URLs, err = store.List(p.Issuer)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		err = t.Execute(w, p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
