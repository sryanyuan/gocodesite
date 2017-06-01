package gocodecc

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type AppConfig struct {
	Debug           bool              `json:"debug"`
	DBAddress       string            `json:"db-address"`
	ListenAddress   string            `json:"listen-address"`
	WeiboAddress    string            `json:"weibo-address"`
	GithubAddress   string            `json:"github-address"`
	CommentProvider string            `json:"comment-provider"`
	CommentContext  map[string]string `json:"comment-context"`
}

var (
	g_appConfig AppConfig
)

func init() {
	g_appConfig.Debug = true
}

// ReadJSONConfig Read config and apply to global config object
func ReadJSONConfig(filename string) error {
	f, err := os.Open(filename)
	if nil != err {
		return err
	}

	fileBytes, err := ioutil.ReadAll(f)
	if nil != err {
		return err
	}

	if err = json.Unmarshal(fileBytes, &g_appConfig); nil != err {
		return err
	}

	return nil
}
