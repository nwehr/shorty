package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
)

func loginCmd(publicUrl string) error {
	loginUrl := publicUrl + "/auth/login"

	fmt.Println("Open the following URL in a browser to continue")
	fmt.Println(loginUrl)

	{
		open := exec.Command("open", loginUrl)
		open.Run()
	}

	done := make(chan bool)

	mux := http.NewServeMux()

	mux.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query()["code"]

		resp, err := http.Get(publicUrl + "/oauth/callback?code=" + code[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		defer resp.Body.Close()

		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		v, err := url.ParseQuery(string(respBytes))
		if err != nil {
			fmt.Println(err)
			return
		}

		accessToken, ok := v["access_token"]
		if !ok {
			fmt.Println("no access token")
			return
		}

		cfg := config{
			Host:        publicUrl,
			AccessToken: accessToken[0],
		}

		err = writeConfig(cfg)
		if err != nil {
			fmt.Println(err)
			return
		}

		w.Write([]byte(`<div>You have successfully logged in. It is safe to close this window.</div>`))

		done <- true
	})

	srv := &http.Server{
		Addr:    "localhost:48162",
		Handler: mux,
	}

	go func() {
		srv.ListenAndServe()
	}()

	<-done

	srv.Shutdown(context.Background())

	fmt.Println("success")

	return nil
}
