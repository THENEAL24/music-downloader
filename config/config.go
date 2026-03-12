package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	HitmoBaseUrl   string `yaml:"base_url"`
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
	MP3Timeout    time.Duration `yaml:"mp3_timeout"`
}

type YtdlpConfig struct {
	Workers        int           `yaml:"workers"`
	RPS            int           `yaml:"rps"`
	Timeout        time.Duration `yaml:"timeout"`
	AudioFormat    string        `yaml:"audio_format"`
	AudioQuality   string        `yaml:"audio_quality"`
	MaxAudioSizeMB int           `yaml:"max_audio_size_mb"`
}

type DownloaderConfig struct {
	Workers          int           `yaml:"workers"`
	HitmoDelay       time.Duration `yaml:"hitmo_delay"`
	UseYtdlpFallback bool          `yaml:"use_ytdlp_fallback"`
}

type ZipConfig struct {
	CompressionLevel int `yaml:"compression_level"`
}

type ServerConfig struct {
	Addr         string        `yaml:"addr"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type APIConfig struct {
	StaticDir     string        `yaml:"static_dir"`
	MaxUploadMB   int64         `yaml:"max_upload_mb"`
	ResultTTL     time.Duration `yaml:"result_ttl"`
	EvictInterval time.Duration `yaml:"evict_interval"`
}

type Config struct {
	App        AppConfig        `yaml:"app"`
	Track      TrackConfig      `yaml:"track"`
	Client     ClientConfig     `yaml:"client"`
	Ytdlp      YtdlpConfig      `yaml:"ytdlp"`
	Downloader DownloaderConfig `yaml:"downloader"`
	Zip        ZipConfig        `yaml:"zip"`
	Server     ServerConfig     `yaml:"server"`
	API        APIConfig        `yaml:"api"`
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
