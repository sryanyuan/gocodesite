package gocodecc

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type AppConfig struct {
	Debug         bool   `json:"debug"`
	DBAddress     string `json:"db-address"`
	ListenAddress string `json:"listen-address"`
	WeiboAddress  string `json:"weibo-address"`
	GithubAddress string `json:"github-address"`
}

var (
	g_appConfig AppConfig
)

func init() {
	g_appConfig.Debug = true
}

// Read config and apply to global config object
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

/*func ReadLuaConfig(filename string, confname string) bool {
	L := luago.LuaGo_newState()
	defer func() {
		L.Destroy()
		L = nil
	}()
	L.OpenStdLibs()
	err := L.LuaGo_SafeDoFile(filename)
	if err != nil {
		seelog.Error(err)
	}

	return L.GetObject(confname, &g_appConfig, true)
}
*/
