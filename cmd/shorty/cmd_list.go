package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func listCmd() error {
	cfg, err := readConfig()
	if err != nil {
		return err
	}

	r, _ := http.NewRequest("GET", cfg.Host+"/api/list", nil)
	r.Header.Add("Authorization", "token "+cfg.AccessToken)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	mappedURLS := []struct {
		ShortURL string `json:"short_url"`
		LongURL  string `json:"long_url"`
	}{}

	json.NewDecoder(resp.Body).Decode(&mappedURLS)

	for _, url := range mappedURLS {
		fmt.Printf("---\n%s\n%s\n\n", url.ShortURL, url.LongURL)
	}

	return nil
}
