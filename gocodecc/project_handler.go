package gocodecc

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var projectCategoryRenderTpls = []string{
	"template/project/category.tpl",
}

var projectArticlesRenderTpls = []string{
	"template/project/articles.tpl",
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
	projectName := vars["projectname"]
	page, err := strconv.Atoi(vars["page"])
	if nil != err {
		ctx.Redirect("/", http.StatusNotFound)
		return
	}
	if page <= 0 {
		page = 1
	}

	articles, pages, err := modelProjectArticleGetArticles(projectName, page-1, 10)
	if nil != err {
		panic(err)
	}

	tplData := make(map[string]interface{})
	tplData["articles"] = articles
	tplData["active"] = "project"
	tplData["project"] = projectName
	tplData["pages"] = pages
	tplData["page"] = page
	data := renderTemplate(ctx, projectArticlesRenderTpls, tplData)
	ctx.w.Write(data)
}

func projectArticleCmdHandler(ctx *RequestContext) {
	vars := mux.Vars(ctx.r)
	cmd := vars["cmd"]

	ctx.w.Write([]byte(cmd))
}
