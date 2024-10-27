package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	OutputDir     string
	Threads       int
	Headers       map[string]string
	ReqType       string
	AuthorChanLen int
	ImgChanLen    int
	RetryChanLen  int
	ProxyURL      string
	ImgQuality    string
	ErrorLogDir   string
}

func LoadConfig(confdir string) (Config, error) {
	bytes, err := os.ReadFile(confdir)
	if err != nil {
		return Config{}, fmt.Errorf("Could not open config file: %w", err)
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("Could not parse config file: %w", err)
	}

	return config, nil
}
