package gocodecc

import (
	"github.com/gorilla/mux"
	"net/http"
)

var memberInfoRenderTpls = []string{
	"template/member/member_info.tpl",
}

func memberInfoHandler(ctx *RequestContext) {
	watchedUser := ctx.user
	vars := mux.Vars(ctx.r)
	member := vars["username"]

	if 0 == ctx.user.Uid ||
		(ctx.user.Uid != 0 && ctx.user.UserName != member) {
		//	is guest, watch other user
		watchedUser = modelWebUserGetUserByUserName(member)
		if nil == watchedUser {
			//	not found
			ctx.Redirect("/", http.StatusFound)
			return
		}
	}

	tplData := make(map[string]interface{})
	tplData["watchedUser"] = watchedUser
	tplData["isSelf"] = (watchedUser.Uid == ctx.user.Uid)
	data := renderTemplate(ctx, memberInfoRenderTpls, tplData)
	ctx.w.Write(data)
}
