package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	HitmoBaseUrl  string `yaml:"base_url"`
	ResultsPerPage int    `yaml:"results_per_page"`
}

type TrackConfig struct {
	ItemsCss  string `yaml:"items_css"`
	ArtistCss string `yaml:"artist_css"`
	TitleCss  string `yaml:"title_css"`
	DlBtnCss  string `yaml:"dl_btn_css"`
}

type ClientConfig struct {
	ClientTimeout time.Duration `yaml:"client_timeout"`
}

type Config struct {
	App AppConfig `yaml:"app"`
	Track TrackConfig `yaml:"track"`
	Client ClientConfig `yaml:"client"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}