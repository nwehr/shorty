package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func createCmd(url string) error {
	cfg, err := readConfig()
	if err != nil {
		return err
	}

	req := struct {
		URL string `json:"long_url"`
	}{
		URL: url,
	}

	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	r, _ := http.NewRequest("POST", cfg.Host+"/api/create", bytes.NewReader(jsonBytes))
	r.Header.Add("Authorization", "token "+cfg.AccessToken)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)

	return nil
}
