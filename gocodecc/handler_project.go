package gocodecc

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"io/ioutil"

	"github.com/cihub/seelog"
	"github.com/dchest/captcha"
	"github.com/gorilla/mux"
)

var projectCategoryRenderTpls = []string{
	"template/project/category.html",
}

var projectArticlesRenderTpls = []string{
	"template/project/articles.html",
	"template/component/article_detail_display.html",
}

var projectArticleNewArticleTpls = []string{
	"template/project/new_article.html",
}

var projectArticleRenderTpls = []string{
	"template/project/article.html",
	"template/component/comment_embed.html",
	"template/component/comment_duoshuo.html",
	"template/component/comment_livere.html",
	"template/component/comment_163.html",
	"template/component/comment_native.html",
	"template/component/reply_list.html",
}

var projectArticleEditArticleRenderTpls = []string{
	"template/project/edit_article.html",
}

func projectHandler(ctx *RequestContext) {
	ctx.Redirect("/project/category", http.StatusFound)
}

func projectCategoryHandler(ctx *RequestContext) {
	//	search all project
	projects, err := modelProjectCategoryGetAll()
	if nil != err {
		panic(err)
	}

	tplData := make(map[string]interface{})
	tplData["category"] = projects
	tplData["active"] = "project"
	data := renderTemplate(ctx, projectCategoryRenderTpls, tplData)
	ctx.w.Write(data)
}

func projectArticlesHandler(ctx *RequestContext) {
	//	search all project
	vars := mux.Vars(ctx.r)
	projectID, err := strconv.Atoi(vars["projectid"])
	if nil != err ||
		0 == projectID {
		ctx.RenderMessagePage("错误", "Invalid projectid", false)
		return
	}
	page, err := strconv.Atoi(vars["page"])
	if nil != err {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}
	if page <= 0 {
		page = 1
	}

	var category ProjectCategoryItem
	err = modelProjectCategoryGetByProjectId(projectID, &category)
	if nil != err {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}

	pageItems := 10
	showPages := 5
	articles, pages, err := modelProjectArticleGetArticles(projectID, page-1, pageItems)
	if nil != err {
		panic(err)
	}

	tplData := make(map[string]interface{})
	tplData["articles"] = articles
	tplData["active"] = "project"
	tplData["project"] = projectID
	tplData["pages"] = pages
	tplData["page"] = page
	tplData["pageItems"] = pageItems
	tplData["showPages"] = showPages
	tplData["category"] = &category
	data := renderTemplate(ctx, projectArticlesRenderTpls, tplData)
	ctx.w.Write(data)
}

func markMessageURLRead(user *WebUser, messageID int, url string) error {
	message, err := modelMessageGetByID(messageID)
	if nil != err {
		return err
	}
	if message.Url != url {
		return errors.New("Url mismatch")
	}
	if message.Receiver != user.Uid {
		return fmt.Errorf("User mismatch receiver %d != %d", message.Receiver, user.Uid)
	}
	return modelMessageMarkRead(user.Uid, messageID)
}

func projectArticleHandler(ctx *RequestContext) {
	vars := mux.Vars(ctx.r)
	articleID, err := strconv.Atoi(vars["articleid"])

	if nil != err {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}

	article, err := modelProjectArticleGet(articleID)
	if nil != err {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}

	// Mark message read?
	ctx.r.ParseForm()
	messageIDStr := ctx.r.Form.Get("messageid")
	if len(messageIDStr) != 0 {
		messageID, err := strconv.Atoi(messageIDStr)
		if nil != err {
			ctx.RenderMessagePage("错误", err.Error(), false)
			return
		}
		if err = markMessageURLRead(ctx.user, messageID, ctx.r.URL.Path); nil != err {
			seelog.Error(err)
		} else {
			modelMessageDelete(messageID)
		}
	}

	//	get author
	author := modelWebUserGetUserByUserName(article.ArticleAuthor)
	if nil == author {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}

	// Increase article visitors
	remoteIP := ctx.GetRemoteIP()
	if remoteIP == "" {
		seelog.Error("Get ip from request failed")
	} else {
		if err = modelArticleVisitorInc(ctx.r.URL.Path, remoteIP); nil != err {
			seelog.Error("Update article visitor failed:", err)
		}
	}

	//	increase click count
	if err = modelProjectArticleIncClick(articleID); nil != err {
		return
	}
	article.Click = article.Click + 1

	// Replies
	var replies []*ReplyModel
	if ctx.config.CommentProvider == "native" {
		// Native should pull all replies
		replies, err = modelReplyGetArticleReply(fmt.Sprintf("/project/%d/article/%d", article.ProjectId, article.Id), 0, 0)
		if nil != err {
			ctx.RenderMessagePage("错误", err.Error(), false)
			return
		}
	}

	tplData := make(map[string]interface{})
	tplData["active"] = "project"
	tplData["article"] = article
	tplData["author"] = author
	tplData["commentID"] = strconv.Itoa(article.Id)
	tplData["commentTitle"] = article.ArticleTitle
	tplData["replies"] = replies
	if ctx.config.CommentProvider == "native" {
		tplData["captchaid"] = captcha.NewLen(4)
	}

	data := renderTemplate(ctx, projectArticleRenderTpls, tplData)
	ctx.w.Write(data)
}

func projectArticleCmdHandler(ctx *RequestContext) {
	vars := mux.Vars(ctx.r)
	cmd := vars["cmd"]
	projectID, err := strconv.Atoi(vars["projectid"])
	cmd = strings.ToLower(cmd)

	if nil != err ||
		0 == projectID {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}

	switch cmd {
	case "new_article":
		{
			_newProjectArticle(ctx, projectID)
		}
	case "edit_article":
		{
			ctx.r.ParseForm()
			articleID, err := strconv.Atoi(ctx.r.Form.Get("articleId"))
			if nil != err {
				ctx.Redirect("/", http.StatusNotFound)
				return
			}

			_editProjectArticle(ctx, articleID)
		}
	default:
		{
			ctx.RenderString("invalid cmd")
		}
	}
}

func _newProjectArticle(ctx *RequestContext, projectId int) {
	//	get category
	var category ProjectCategoryItem
	err := modelProjectCategoryGetByProjectId(projectId, &category)
	if nil != err {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}

	//	check auth
	if ctx.user.Uid == 0 {
		ctx.RenderString("acess denied")
		return
	}
	if ctx.user.NickName == category.Author ||
		ctx.user.Permission >= category.PostPriv {
		//	nothing
	} else {
		ctx.RenderString("access denied")
		return
	}

	tplData := make(map[string]interface{})
	tplData["active"] = "project"
	tplData["project"] = &category
	tplData["captchaid"] = captcha.NewLen(4)
	data := renderTemplate(ctx, projectArticleNewArticleTpls, tplData)
	ctx.w.Write(data)
}

func _editProjectArticle(ctx *RequestContext, articleId int) {
	article, err := modelProjectArticleGet(articleId)
	if err != nil {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}

	tplData := make(map[string]interface{})
	tplData["active"] = "project"
	tplData["article"] = article
	tplData["captchaid"] = captcha.NewLen(4)
	data := renderTemplate(ctx, projectArticleEditArticleRenderTpls, tplData)
	ctx.w.Write(data)
}

func projectArticleReplyHandler(ctx *RequestContext) {
	if ctx.r.Method != http.MethodGet {
		ctx.RenderMessagePage("Permission denied", "Invalid method", false)
		return
	}
	vars := mux.Vars(ctx.r)
	articleID, err := strconv.Atoi(vars["articleid"])

	if nil != err {
		seelog.Debug(err)
		ctx.Redirect("/", http.StatusNotFound)
		return
	}

	// Read reply body
	bodyBytes, err := ioutil.ReadAll(ctx.r.Body)
	if nil != err {
		seelog.Debug(err)
		ctx.WriteResponse([]byte(err.Error()))
		return
	}

	// Is use logined ?
	if ctx.user.Uid == 0 {
		seelog.Debug("sign in")
		ctx.Redirect("/signin", http.StatusContinue)
		return
	}

	seelog.Debug(string(bodyBytes), articleID)
}
