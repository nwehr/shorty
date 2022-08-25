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

	r, _ := http.NewRequest("GET", cfg.Host+"/shorturls", nil)
	r.Header.Add("Authorization", "token "+cfg.AccessToken)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	shortUrls := []struct {
		ID  string
		URL string
	}{}

	json.NewDecoder(resp.Body).Decode(&shortUrls)

	for _, shortUrl := range shortUrls {
		fmt.Printf("---\n%s/%s\n%s\n\n", cfg.Host, shortUrl.ID, shortUrl.URL)
	}

	return nil
}
