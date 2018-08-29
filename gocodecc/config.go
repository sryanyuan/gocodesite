package gocodecc

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/go-validator/validator"
)

type CDNConfig struct {
	JQueryJS            string `json:"jquery-js" toml:"jquery-js"`
	JQueryCSS           string `json:"jquery-css" toml:"jquery-css"`
	JQueryPlaceholderJS string `json:"jquery-placeholder-js" toml:"jquery-placeholder-js"`
	BootstrapJS         string `json:"bootstrap-js" toml:"bootstrap-js"`
	BootstrapCSS        string `json:"bootstrap-css" toml:"bootstrap-css"`
	BootstrapThemeCSS   string `json:"bootstrap-theme-css" toml:"bootstrap-theme-css"`
	FontAwesomeCSS      string `json:"font-awesome-css" toml:"font-awesome-css"`
	JQueryQRCodeJS      string `json:"jquery-qrcode-js" toml:"jquery-qrcode-js"`
}

type AppConfig struct {
	Debug           bool              `json:"debug" toml:"debug"`
	DBAddress       string            `json:"db-address" toml:"db-address" validate:"nonzero"`
	SiteTitle       string            `json:"site-title" toml:"site-title"`
	HomeTitle       string            `json:"home-title" toml:"home-title"`
	BannerImage     string            `json:"banner-image" toml:"banner-image"`
	BannerText      string            `json:"banner-text" toml:"banner-text"`
	FooterCopyright string            `json:"footer-copyright" toml:"footer-copyright"`
	AboutHTMLFile   string            `json:"about-html-file" toml:"about-html-file"`
	ResumeFile      string            `json:"resume-file" toml:"resume-file"`
	ListenAddress   string            `json:"listen-address" toml:"listen-address" validate:"nonzero"`
	WeiboAddress    string            `json:"weibo-address" toml:"weibo-address"`
	GithubAddress   string            `json:"github-address" toml:"github-address"`
	NginxProxy      bool              `json:"nginx-proxy" toml:"nginx-proxy"`
	Timezone        string            `json:"timezone" toml:"timezone"`
	CommentProvider string            `json:"comment-provider" toml:"comment-provider"`
	CommentContext  map[string]string `json:"comment-context" toml:"comment-context"`
	CallSecret      string            `json:"call-secret" toml:"call-secret"`
	DonateCall      string            `json:"donate-call" toml:"donate-call"`
	MsgPush         MsgPushConfig     `toml:"msg-push"`
	CDN             CDNConfig         `json:"cdn" toml:"cdn"`
	Domain          string            `json:"domain" toml:"domain"`
	NewBlog         string            `json:"new-blog" toml:"new-blog"`
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
