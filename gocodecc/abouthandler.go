package gocodecc

import (
	"html/template"
	"net/http"

	"github.com/cihub/seelog"
)

func aboutHander(ctx *RequestContext) {
	t, err := template.ParseFiles("template/layout.tpl")
	if nil != err {
		seelog.Error(err)
		http.Error(ctx.w, ctx.r, http.StatusInternalServerError)
		return
	}
}
