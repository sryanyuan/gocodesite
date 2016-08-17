package gocodecc

import (
	"net/http"
	"strconv"

	"github.com/cihub/seelog"
	"github.com/gorilla/mux"
)

var layoutFiles = []string{
	"template/layout.tpl",
	"template/component/navbar_v2.tpl",
	"template/component/footer.tpl",
}

var (
	metaInfoCreateSiteTime int64
)

const kErrMsg_InternalError = "内部错误，请重试"

func Start() error {
	var err error
	//	initialize model
	if err = initModels(); nil != err {
		seelog.Error(err)
	}

	//	initialize meta info
	initMetaInfo()

	//	get base meta info
	metaInfoCreateSiteTimeStr, err := modelMetaInfoGet("create_site_time")
	if nil != err {
		seelog.Error("Failed to read meta info")
		return err
	}
	metaInfoCreateSiteTime, _ = strconv.ParseInt(metaInfoCreateSiteTimeStr, 10, 64)

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
