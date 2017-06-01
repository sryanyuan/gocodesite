package gocodecc

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/go-validator/validator"
)

type AppConfig struct {
	Debug           bool              `json:"debug" toml:"debug"`
	DBAddress       string            `json:"db-address" toml:"db-address" validate:"nonzero"`
	SiteTitle       string            `json:"site-title" toml:"site-title"`
	HomeTitle       string            `json:"home-title" toml:"home-title"`
	ListenAddress   string            `json:"listen-address" toml:"listen-address" validate:"nonzero"`
	WeiboAddress    string            `json:"weibo-address" toml:"weibo-address"`
	GithubAddress   string            `json:"github-address" toml:"github-address"`
	CommentProvider string            `json:"comment-provider" toml:"comment-provider"`
	CommentContext  map[string]string `json:"comment-context" toml:"comment-context"`
}

// ReadJSONConfig returns config object loading from json format config file
func ReadJSONConfig(filename string) (*AppConfig, error) {
	f, err := os.Open(filename)
	if nil != err {
		return nil, err
	}

	fileBytes, err := ioutil.ReadAll(f)
	if nil != err {
		return nil, err
	}

	var config AppConfig
	config.Debug = true
	if err = json.Unmarshal(fileBytes, &config); nil != err {
		return nil, err
	}

	return &config, nil
}

// ReadTOMLConfig returns config object loading from toml format config file
func ReadTOMLConfig(filename string) (*AppConfig, error) {
	var config AppConfig
	config.Debug = true
	if _, err := toml.DecodeFile(filename, &config); nil != err {
		return nil, err
	}
	if err := validator.Validate(&config); nil != err {
		return nil, err
	}
	return &config, nil
}
