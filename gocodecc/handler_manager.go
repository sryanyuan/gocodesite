package gocodecc

import (
	"net/http"

	"github.com/gorilla/mux"
)

var managerIndexRenderTpls = []string{
	"template/manager/layout.html",
	"template/manager/leftmenu.html",
	"template/manager/users.html",
	"template/manager/bmkv.html",
}

func managerPanelHandler(ctx *RequestContext) {
	vars := mux.Vars(ctx.r)
	panel := vars["panel"]

	// Dispatch to each panel
	switch panel {
	case "users":
		{
			managerUserHandler(ctx)
		}
	case "bmkv":
		{
			managerBmkvHandler(ctx)
		}
	default:
		{
			managerDefaultHandler(ctx)
		}
	}
}

func managerHandler(ctx *RequestContext) {
	ctx.Redirect("/manager/users", http.StatusFound)
}

func managerUserHandler(ctx *RequestContext) {
	users, err := modelWebUserGetAll(0, 0)
	if nil != err {
		ctx.RenderMessagePage("错误", err.Error(), false)
		return
	}

	tplData := make(map[string]interface{})
	tplData["active"] = "users"
	tplData["users"] = users
	data := renderTemplate(ctx, managerIndexRenderTpls, tplData)
	ctx.w.Write(data)
}

func managerBmkvHandler(ctx *RequestContext) {
	kvs, err := modelBmkvGetAll(0, 0)
	if nil != err {
		ctx.RenderMessagePage("错误", err.Error(), false)
		return
	}

	tplData := make(map[string]interface{})
	tplData["active"] = "bmkv"
	tplData["kvs"] = kvs
	data := renderTemplate(ctx, managerIndexRenderTpls, tplData)
	ctx.w.Write(data)
}

func managerDefaultHandler(ctx *RequestContext) {
	data := renderTemplate(ctx, managerIndexRenderTpls, nil)
	ctx.w.Write(data)
}
