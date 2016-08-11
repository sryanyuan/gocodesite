package gocodecc

import (
	"fmt"
	"time"

	"github.com/gorilla/mux"
)

type AjaxResult struct {
	Result int    `json:"Result"`
	Msg    string `json:"Msg"`
}

func ajaxHandler(ctx *RequestContext) {
	vars := mux.Vars(ctx.r)
	action := vars["action"]
	var result AjaxResult
	result.Result = -1
	defer renderJson(ctx, &result)

	switch action {
	case "project_create":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			//	must be superadmin
			if ctx.user.Permission < kPermission_SuperAdmin {
				result.Msg = "permission denied"
				return
			}

			//	check project name and project describe
			ctx.r.ParseForm()
			projectName := ctx.r.Form.Get("project[name]")
			projectDescribe := ctx.r.Form.Get("project[describe]")
			projectImage := ctx.r.Form.Get("project[image]")
			ctx.r.Body.Close()

			if len(projectName) == 0 ||
				len(projectDescribe) == 0 {
				result.Msg = "invalid project name or project describe"
				return
			}

			var project ProjectCategoryItem
			project.Author = ctx.user.NickName
			project.Image = projectImage
			project.ProjectName = projectName
			project.ProjectDescribe = projectDescribe
			err := modelProjectCategoryAdd(&project)
			if nil != err {
				result.Msg = err.Error()
				return
			}

			//	done
			result.Result = 0
		}
	case "project_edit":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			//	must be superadmin
			if ctx.user.Permission < kPermission_SuperAdmin {
				result.Msg = "permission denied"
				return
			}

			//	check project name and project describe
			ctx.r.ParseForm()
			var err error
			projectName := ctx.r.Form.Get("project[name]")
			projectDescribe := ctx.r.Form.Get("project[describe]")
			projectImage := ctx.r.Form.Get("project[image]")
			ctx.r.Body.Close()

			if len(projectName) == 0 ||
				len(projectDescribe) == 0 {
				result.Msg = "invalid project name or project describe"
				return
			}

			//	get the original item
			var originPrj ProjectCategoryItem
			if err := modelProjectCategoryGetByProjectName(projectName, &originPrj); nil != err {
				result.Msg = err.Error()
				return
			}

			originPrj.ProjectName = projectName
			originPrj.ProjectDescribe = projectDescribe
			originPrj.Image = projectImage
			err = modelProjectCategoryUpdateProject(&originPrj)
			if nil != err {
				result.Msg = err.Error()
				return
			}

			//	done
			result.Result = 0
		}
	case "project_delete":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			//	must be superadmin
			if ctx.user.Permission < kPermission_SuperAdmin {
				result.Msg = "permission denied"
				return
			}

			ctx.r.ParseForm()
			projectName := ctx.r.Form.Get("project[name]")
			ctx.r.Body.Close()

			if len(projectName) == 0 {
				result.Msg = "invalid project name"
				return
			}

			err := modelProjectCategoryRemoveByProjectName(projectName)
			if nil != err {
				result.Msg = err.Error()
				return
			}

			//	done
			result.Result = 0
		}
	case "article_submit":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			ctx.r.ParseForm()
			projectName := ctx.r.Form.Get("project")
			if len(projectName) == 0 {
				result.Msg = "invalid project"
				return
			}
			ctx.r.Body.Close()
			//	check auth
			var prj ProjectCategoryItem
			if err := modelProjectCategoryGetByProjectName(projectName, &prj); nil != err {
				result.Msg = err.Error()
				return
			}
			//	check auth
			if ctx.user.Permission < prj.PostPriv {
				result.Msg = "permission denied"
				return
			}
			//	check valid
			title := ctx.r.Form.Get("title")
			if len(title) >= 128 {
				result.Msg = "标题长度太长了"
				return
			}
			if len(title) == 0 {
				result.Msg = "请输入标题"
				return
			}
			contentHtml := ctx.r.Form.Get("editormd-html-code")
			if len(contentHtml) == 0 {
				result.Msg = "请输入内容"
				return
			}
			if len(contentHtml) >= 12800 {
				result.Msg = "内容太长了"
				return
			}
			contentMarkdown := ctx.r.Form.Get("editormd-markdown-doc")
			if len(contentMarkdown) == 0 {
				result.Msg = "请输入内容"
				return
			}
			if len(contentMarkdown) >= 12800 {
				result.Msg = "内容太长了"
				return
			}
			//	do post
			var postArticle ProjectArticleItem
			postArticle.ActiveTime = time.Now().Unix()
			postArticle.PostTime = time.Now().Unix()
			postArticle.ArticleTitle = title
			postArticle.ArticleAuthor = ctx.user.NickName
			postArticle.ArticleContentHtml = contentHtml
			postArticle.ArticleContentMarkdown = contentMarkdown
			postArticle.ProjectName = projectName
			articleId, err := modelProjectArticleNewArticle(&postArticle)
			if nil != err {
				result.Msg = err.Error()
				return
			}
			//	done
			result.Result = 0
			result.Msg = fmt.Sprintf("/project/%s/article/%d", projectName, articleId)
		}
	default:
		{
			result.Msg = "invalid ajax request"
		}
	}
}
