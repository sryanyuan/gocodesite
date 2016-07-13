package gocodecc

import (
	"net/http"

	"github.com/gorilla/mux"
)

var projectCategoryRenderTpls = []string{
	"template/project/category.tpl",
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

func projectContentHandler(ctx *RequestContext) {
	//	search all project
	vars := mux.Vars(ctx.r)
	member := vars["projectname"]
	ctx.WriteResponse([]byte(member))
}
