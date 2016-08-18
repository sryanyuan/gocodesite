package gocodecc

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var memberInfoRenderTpls = []string{
	"template/member/member_info.tpl",
}

var memberArticlesRenderTpls = []string{
	"template/member/member_articles.tpl",
	"template/component/article_detail_display.tpl",
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

	socialInfo, _ := modelSocialInfoGet(watchedUser.Uid)
	//socialInfo.Weibo = "hh"
	//socialInfo.Github = "hh"

	//	get articles
	articles, err := modelProjectArticleGetByAuthor(watchedUser.NickName, 0, 10)
	if nil != err {
		ctx.RenderMessagePage("错误", err.Error(), false)
		return
	}

	//	get article count
	articleCount, err := modelProjectArticleGetArticleCountByAuthor(watchedUser.NickName)
	if nil != err {
		ctx.RenderMessagePage("错误", err.Error(), false)
		return
	}

	tplData := make(map[string]interface{})
	tplData["watchedUser"] = watchedUser
	tplData["isSelf"] = (watchedUser.Uid == ctx.user.Uid)
	tplData["replyCount"] = 0
	tplData["postCount"] = articleCount
	tplData["watchedSocialInfo"] = socialInfo
	tplData["articles"] = articles
	data := renderTemplate(ctx, memberInfoRenderTpls, tplData)
	ctx.w.Write(data)
}

func memberArticlesHandler(ctx *RequestContext) {
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

	ctx.r.ParseForm()
	pageStr := ctx.r.Form.Get("p")
	var page int
	var err error
	if len(pageStr) == 0 {
		page = 1
	} else {
		page, err = strconv.Atoi(pageStr)
		if nil != err {
			ctx.RenderMessagePage("错误", err.Error(), false)
			return
		}
	}

	//	get total page
	articleCount, err := modelProjectArticleGetArticleCountByAuthor(member)
	if nil != err {
		ctx.RenderMessagePage("错误", err.Error(), false)
		return
	}

	articlePerPage := 10
	showPages := 5
	pages := (articleCount + articlePerPage - 1) / articlePerPage
	var articles []*ProjectArticleItem

	if 0 == pages {
		//	never post
		if page != 1 {
			ctx.RenderMessagePage("错误", kErrMsg_InternalError, false)
			return
		}
		articles = make([]*ProjectArticleItem, 0, 1)
	} else {
		if page <= 0 ||
			page > pages {
			ctx.RenderMessagePage("错误", kErrMsg_InternalError, false)
			return
		}

		//	get articles
		articles, err = modelProjectArticleGetByAuthor(watchedUser.NickName, page-1, articlePerPage)
		if nil != err {
			ctx.RenderMessagePage("错误", kErrMsg_InternalError, false)
			return
		}
	}

	tplData := make(map[string]interface{})
	tplData["watchedUser"] = watchedUser
	tplData["isSelf"] = (watchedUser.Uid == ctx.user.Uid)
	tplData["articles"] = articles
	tplData["pages"] = pages
	tplData["page"] = page
	tplData["showPages"] = showPages
	data := renderTemplate(ctx, memberArticlesRenderTpls, tplData)
	ctx.w.Write(data)
}
