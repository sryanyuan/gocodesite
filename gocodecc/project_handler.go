package gocodecc

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

var projectCategoryRenderTpls = []string{
	"template/project/category.tpl",
}

var projectArticlesRenderTpls = []string{
	"template/project/articles.tpl",
	"template/component/article_detail_display.tpl",
}

var projectArticleNewArticleTpls = []string{
	"template/project/new_article.tpl",
}

var projectArticleRenderTpls = []string{
	"template/project/article.tpl",
	"template/component/comment_article_html.tpl",
}

var projectArticleEditArticleRenderTpls = []string{
	"template/project/edit_article.tpl",
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
	projectId, err := strconv.Atoi(vars["projectid"])
	if nil != err ||
		0 == projectId {
		ctx.Redirect("/", http.StatusNotFound)
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
	err = modelProjectCategoryGetByProjectId(projectId, &category)
	if nil != err {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}

	pageItems := 10
	showPages := 5
	articles, pages, err := modelProjectArticleGetArticles(projectId, page-1, pageItems)
	if nil != err {
		panic(err)
	}

	tplData := make(map[string]interface{})
	tplData["articles"] = articles
	tplData["active"] = "project"
	tplData["project"] = projectId
	tplData["pages"] = pages
	tplData["page"] = page
	tplData["pageItems"] = pageItems
	tplData["showPages"] = showPages
	tplData["category"] = &category
	data := renderTemplate(ctx, projectArticlesRenderTpls, tplData)
	ctx.w.Write(data)
}

func projectArticleHandler(ctx *RequestContext) {
	vars := mux.Vars(ctx.r)
	articleId, err := strconv.Atoi(vars["articleid"])

	if nil != err {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}

	article, err := modelProjectArticleGet(articleId)
	if nil != err {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}

	//	get author
	author := modelWebUserGetUserByUserName(article.ArticleAuthor)
	if nil == author {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}

	//	increase click count
	if err = modelProjectArticleIncClick(articleId); nil != err {
		return
	}
	article.Click = article.Click + 1

	tplData := make(map[string]interface{})
	tplData["active"] = "project"
	tplData["article"] = article
	tplData["author"] = author
	data := renderTemplate(ctx, projectArticleRenderTpls, tplData)
	ctx.w.Write(data)
}

func projectArticleCmdHandler(ctx *RequestContext) {
	vars := mux.Vars(ctx.r)
	cmd := vars["cmd"]
	projectId, err := strconv.Atoi(vars["projectid"])
	cmd = strings.ToLower(cmd)

	if nil != err ||
		0 == projectId {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}

	switch cmd {
	case "new_article":
		{
			_newProjectArticle(ctx, projectId)
		}
	case "edit_article":
		{
			ctx.r.ParseForm()
			articleId, err := strconv.Atoi(ctx.r.Form.Get("articleId"))
			if nil != err {
				ctx.Redirect("/", http.StatusNotFound)
				return
			}

			_editProjectArticle(ctx, articleId)
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
	data := renderTemplate(ctx, projectArticleEditArticleRenderTpls, tplData)
	ctx.w.Write(data)
}
