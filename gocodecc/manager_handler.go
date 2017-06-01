package gocodecc

import "github.com/cihub/seelog"

var managerIndexRenderTpls = []string{
	"template/manager/index.html",
	"template/manager/leftmenu.html",
}

func managerHandler(ctx *RequestContext) {
	managerUserHandler(ctx)
}

func managerUserHandler(ctx *RequestContext) {
	users, err := modelWebUserGetAll(0, 0)
	if nil != err {
		ctx.RenderMessagePage("错误", err.Error(), false)
		return
	}

	tplData := make(map[string]interface{})
	tplData["active"] = "user"
	tplData["users"] = users
	seelog.Debug(len(users))
	data := renderTemplate(ctx, managerIndexRenderTpls, tplData)
	ctx.w.Write(data)
}
