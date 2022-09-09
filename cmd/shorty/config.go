package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	Host        string `yaml:"host"`
	AccessToken string `yaml:"access_token"`
}

func readConfig() (config, error) {
	cfg := config{}

	f, err := os.Open(getConfigPath())
	if err != nil {
		return cfg, err
	}

	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func writeConfig(cfg config) error {
	f, err := os.OpenFile(getConfigPath(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	err = yaml.NewEncoder(f).Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}

func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	configDir := homeDir + "/.config"

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.Mkdir(configDir, os.ModeDir)
	}

	return configDir + "/shorty.yml"
}
