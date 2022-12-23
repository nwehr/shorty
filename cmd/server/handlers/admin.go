package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/nwehr/shorty/cmd/server/models"
	"github.com/nwehr/shorty/cmd/server/options"
)

type (
	PageData struct {
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
		t, err := template.ParseFiles("templates/index.tmpl", "templates/header.tmpl", "templates/footer.tmpl", "templates/navbar.tmpl", "templates/content.tmpl")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		authorization, err := getAuthorization(r)
		if err != nil {
			fmt.Println(err)
		}

		data := PageData{
			Opts:     opts,
			LoggedIn: err == nil,
			Issuer:   AccessToken(authorization.AccessToken).GetIssuer(),
			Username: AccessToken(authorization.AccessToken).GetUsername(),
		}

		if err == nil {
			data.URLs, err = store.List(data.Issuer)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		err = t.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
