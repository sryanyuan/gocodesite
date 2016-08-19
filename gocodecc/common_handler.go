package gocodecc

import (
	"github.com/gorilla/mux"
)

var commonMessageRenderTpls []string = []string{
	"template/common/message.tpl",
}

var commonDownloadRenderTpls []string = []string{
	"template/common/download.tpl",
}

func commonHandler(ctx *RequestContext) {
	vars := mux.Vars(ctx.r)
	action := vars["action"]

	switch action {
	case "message":
		{
			ctx.r.ParseForm()
			tplData := make(map[string]interface{})
			tplData["Text"] = ctx.r.Form.Get("text")
			tplData["Title"] = ctx.r.Form.Get("title")
			tplData["Result"] = ctx.r.Form.Get("result")
			data := renderTemplate(ctx, commonMessageRenderTpls, tplData)
			ctx.w.Write(data)
		}
	case "download":
		{
			ctx.r.ParseForm()
			tplData := make(map[string]interface{})
			tplData["Text"] = ctx.r.Form.Get("text")
			tplData["Title"] = ctx.r.Form.Get("title")
			tplData["Url"] = ctx.r.Form.Get("url")
			data := renderTemplate(ctx, commonDownloadRenderTpls, tplData)
			ctx.w.Write(data)
		}
	}
}
