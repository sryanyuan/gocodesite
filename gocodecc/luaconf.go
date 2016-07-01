package gocodecc

import (
	"github.com/cihub/seelog"
	"github.com/sryanyuan/luago"
)

type AppConfig struct {
	Debug         bool
	DBAddress     string
	ListenAddress string
	WeiboAddress  string
	GithubAddress string
}

var (
	g_appConfig AppConfig
)

func init() {
	g_appConfig.Debug = true
}

func ReadLuaConfig(filename string, confname string) bool {
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
