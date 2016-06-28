package gocodecc

import (
	"errors"

	"net/http"

	"github.com/cihub/seelog"
	"github.com/gorilla/mux"
)

var layoutFiles = []string{
	"template/layout.tpl",
	"template/component/navbar.tpl",
	"template/component/footer.tpl",
}

func Start() error {
	var err error
	//	check db config
	if len(g_appConfig.DBAddress) == 0 {
		seelog.Error("Invalid config")
		return errors.New("Invalid config")
	}

	//	initialize routers
	r := mux.NewRouter()
	InitRouters(r)
	http.Handle("/", r)

	//	run the server
	seelog.Info("Http server listen on:", g_appConfig.ListenAddress)
	err = http.ListenAndServe(g_appConfig.ListenAddress, nil)
	if nil != err {
		seelog.Error("Http error:", err)
		return err
	}

	return nil
}
