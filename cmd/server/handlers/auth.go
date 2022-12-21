package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nwehr/shorty/cmd/server/options"
)

type (
	AuthHandler struct {
		Opts options.Options
	}

	Authorization struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token" `
	}

	AccessToken  string
	RefreshToken string
)

func (t AccessToken) GetIssuer() string {
	segments := strings.Split(string(t), ".")

	if len(segments) != 3 {
		fmt.Println("len(segments): ", len(segments))
		return ""
	}

	decoded, err := base64.RawStdEncoding.DecodeString(segments[1])
	if err != nil {
		return ""
	}

	s := struct {
		Issuer string `json:"iss"`
	}{}

	err = json.Unmarshal(decoded, &s)
	if err != nil {
		return ""
	}

	return s.Issuer
}

func (t AccessToken) GetUsername() string {
	segments := strings.Split(string(t), ".")

	if len(segments) != 3 {
		return ""
	}

	decoded, err := base64.RawStdEncoding.DecodeString(segments[1])
	if err != nil {
		return ""
	}

	s := struct {
		Username string `json:"preferred_username"`
	}{}

	err = json.Unmarshal(decoded, &s)
	if err != nil {
		return ""
	}

	return s.Username
}

func (s AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqeustsCounter.Inc()

	if r.Method == "GET" {
		switch r.URL.Path {
		case "/auth/callback":
			s.handleCallback(w, r)
			return
		case "/auth/login":
			s.handleLogin(w, r)
			return
		case "/auth/logout":
			s.handleLogout(w, r)
		}
	}

	http.Error(w, "", http.StatusBadRequest)
}

func (s AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if s.Opts.ClientID == "" {
		http.Error(w, "no oauth client id", http.StatusInternalServerError)
		return
	}

	v := url.Values{}
	v.Add("client_id", s.Opts.ClientID)
	v.Add("redirect_uri", s.Opts.RedirectURL)
	v.Add("response_type", "code")

	authURL, _ := url.Parse(s.Opts.AuthorizationURL)
	authURL.RawQuery = v.Encode()

	http.Redirect(w, r, authURL.String(), http.StatusTemporaryRedirect)
}

func (s AuthHandler) handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:    "AUTHORIZATION",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}

	http.SetCookie(w, cookie)
	http.Redirect(w, r, s.Opts.LogoutURL+"?redirect_uri="+s.Opts.PublicURL+"/admin", http.StatusTemporaryRedirect)
}

func (s AuthHandler) handleCallback(w http.ResponseWriter, r *http.Request) {
	form := url.Values{}
	form.Set("client_id", s.Opts.ClientID)
	form.Set("client_secret", s.Opts.ClientSecret)
	form.Set("redirect_uri", s.Opts.RedirectURL)
	form.Set("grant_type", "authorization_code")
	form.Set("code", r.URL.Query()["code"][0])

	resp, err := http.DefaultClient.PostForm(s.Opts.TokenURL, form)
	if err != nil {
		s.Opts.Logger.Errorf("could not handle callback: %v\n", err)
		errorCounter.Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	authorization := Authorization{}

	err = json.NewDecoder(resp.Body).Decode(&authorization)
	if err != nil {
		s.Opts.Logger.Errorf("could not decode callback response: %v\n", err)
		return
	}

	if err != nil {
		s.Opts.Logger.Errorf("could not decode callback response: %v\n", err)
		errorCounter.Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authBytes, _ := json.Marshal(authorization)

	cookie := http.Cookie{
		Name:  "AUTHORIZATION",
		Value: base64.RawStdEncoding.EncodeToString(authBytes),
		Path:  "/",
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, s.Opts.PublicURL+"/admin", http.StatusTemporaryRedirect)
}

func getAuthorization(r *http.Request) (Authorization, error) {
	cookie, err := r.Cookie("AUTHORIZATION")
	if err != nil {
		return Authorization{}, fmt.Errorf("could not get cookie: %w", err)
	}

	authorization := Authorization{}
	decoded, _ := base64.RawStdEncoding.DecodeString(cookie.Value)

	err = json.Unmarshal(decoded, &authorization)
	if err != nil {
		return authorization, fmt.Errorf("could not parse authorization: %w", err)
	}

	return authorization, nil
}
