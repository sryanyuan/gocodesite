package gocodecc

import (
	"fmt"
	"strconv"
	"strings"
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
				len(projectDescribe) == 0 ||
				len(projectName) >= kCategoryNameLimit ||
				len(projectDescribe) >= kCategoryDescribeLimit {
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
			projectOldName := ctx.r.Form.Get("project[oldname]")
			projectName := ctx.r.Form.Get("project[name]")
			projectDescribe := ctx.r.Form.Get("project[describe]")
			projectImage := ctx.r.Form.Get("project[image]")
			ctx.r.Body.Close()

			if len(projectName) == 0 ||
				len(projectDescribe) == 0 ||
				len(projectOldName) == 0 ||
				len(projectName) >= kCategoryNameLimit ||
				len(projectDescribe) >= kCategoryDescribeLimit {
				result.Msg = "invalid project name or project describe"
				return
			}

			//	get the original item
			var originPrj ProjectCategoryItem
			if err := modelProjectCategoryGetByProjectName(projectOldName, &originPrj); nil != err {
				result.Msg = err.Error()
				return
			}

			if originPrj.ProjectName == projectName &&
				originPrj.ProjectDescribe == projectDescribe &&
				originPrj.Image == projectImage {
				return
			}

			var newPrj ProjectCategoryItem
			newPrj = originPrj
			newPrj.ProjectName = projectName
			newPrj.ProjectDescribe = projectDescribe
			newPrj.Image = projectImage
			err = modelProjectCategoryUpdateProject(&originPrj, &newPrj)
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
			if len(title) >= kArticleTitleLimit {
				result.Msg = "标题长度太长了"
				return
			}
			if len(title) == 0 {
				result.Msg = "请输入标题"
				return
			}
			contentHtml := ctx.r.Form.Get("editormd-html-code")
			contentHtml = strings.Replace(contentHtml, "<pre>", "", -1)
			contentHtml = strings.Replace(contentHtml, "</pre>", "", -1)
			if len(contentHtml) == 0 {
				result.Msg = "请输入内容"
				return
			}
			if len(contentHtml) >= kArticleContentLimit {
				result.Msg = "内容太长了"
				return
			}
			contentMarkdown := ctx.r.Form.Get("editormd-markdown-doc")
			if len(contentMarkdown) == 0 {
				result.Msg = "请输入内容"
				return
			}
			if len(contentMarkdown) >= kArticleContentLimit {
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
			postArticle.ProjectId = prj.Id
			articleId, err := modelProjectArticleNewArticle(&postArticle)
			if nil != err {
				result.Msg = err.Error()
				return
			}
			//	done
			result.Result = 0
			result.Msg = fmt.Sprintf("/project/%s/article/%d", projectName, articleId)
		}
	case "article_edit":
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
			articleId, _ := strconv.Atoi(ctx.r.Form.Get("articleId"))
			if 0 == articleId {
				result.Msg = "invalid articleId"
				return
			}
			ctx.r.Body.Close()
			//	check auth
			article, err := modelProjectArticleGet(articleId)
			if nil != err {
				result.Msg = err.Error()
				return
			}
			if article.ArticleAuthor != ctx.user.NickName {
				if ctx.user.Permission < kPermission_SuperAdmin {
					result.Msg = "access denied"
					return
				}
			}
			//	check valid
			title := ctx.r.Form.Get("title")
			if len(title) >= kArticleTitleLimit {
				result.Msg = "标题长度太长了"
				return
			}
			if len(title) == 0 {
				result.Msg = "请输入标题"
				return
			}
			contentHtml := ctx.r.Form.Get("editormd-html-code")
			contentHtml = strings.Replace(contentHtml, "<pre>", "", -1)
			contentHtml = strings.Replace(contentHtml, "</pre>", "", -1)
			if len(contentHtml) == 0 {
				result.Msg = "请输入内容"
				return
			}
			if len(contentHtml) >= kArticleContentLimit {
				result.Msg = "内容太长了"
				return
			}
			contentMarkdown := ctx.r.Form.Get("editormd-markdown-doc")
			if len(contentMarkdown) == 0 {
				result.Msg = "请输入内容"
				return
			}
			if len(contentMarkdown) >= kArticleContentLimit {
				result.Msg = "内容太长了"
				return
			}
			//	do post
			colsEdit := []string{"active_time", "edit_time"}
			article.ActiveTime = time.Now().Unix()
			article.EditTime = time.Now().Unix()
			if article.ArticleTitle != title {
				article.ArticleTitle = title
				colsEdit = append(colsEdit, "article_title")
			}
			if article.ArticleContentHtml != contentHtml {
				article.ArticleContentHtml = contentHtml
				article.ArticleContentMarkdown = contentMarkdown
				colsEdit = append(colsEdit, "article_content_html")
				colsEdit = append(colsEdit, "article_content_markdown")
			}
			_, err = modelProjectArticleEditArticle(article, colsEdit)
			if nil != err {
				result.Msg = err.Error()
				return
			}
			//	done
			result.Result = 0
			result.Msg = fmt.Sprintf("/project/%s/article/%d", projectName, articleId)
		}
	case "article_delete":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			ctx.r.ParseForm()
			articleId, err := strconv.Atoi(ctx.r.Form.Get("articleId"))
			ctx.r.Body.Close()
			if err != nil ||
				0 == articleId {
				result.Msg = "invalid articleId"
				return
			}

			//	get article
			article, err := modelProjectArticleGet(articleId)
			if nil != err {
				result.Msg = "invalid article"
				return
			}

			//	must be superadmin
			if ctx.user.Permission <= kPermission_Admin {
				result.Msg = "access denied"
				return
			}

			err = modelProjectArticleDelete(articleId, article.ProjectId)
			if nil != err {
				result.Msg = "delete article failed"
				return
			}

			//	done
			result.Result = 0
			result.Msg = fmt.Sprintf("/project/%s/page/1", article.ProjectName)
		}
	case "article_top":
		{
			if ctx.r.Method != "POST" {
				result.Msg = "invalid method"
				return
			}

			ctx.r.ParseForm()
			fmt.Println(ctx.r.Form)
			defer ctx.r.Body.Close()
			articleId, err := strconv.Atoi(ctx.r.Form.Get("articleId"))
			if err != nil ||
				0 == articleId {
				result.Msg = "invalid articleId"
				return
			}
			top, err := strconv.Atoi(ctx.r.Form.Get("top"))
			if err != nil {
				result.Msg = "invalid top"
				return
			}

			doTop := true
			if 0 == top {
				doTop = false
			}

			err = modelProjectArticleSetTop(articleId, doTop)
			if nil != err {
				result.Msg = "set top failed"
				return
			}

			//	done
			result.Result = 0
		}
	case "account_verify":
		{
			if ctx.r.Method != "GET" {
				result.Msg = "invalid method"
				return
			}

			ctx.r.ParseForm()
			account := ctx.r.Form.Get("account")
			password := ctx.r.Form.Get("password")

			if len(account) == 0 ||
				len(password) == 0 ||
				len(account) > 20 ||
				len(password) > 100 {
				result.Msg = "invalid input"
				return
			}

			user := modelWebUserGetUserByUserName(account)
			if nil == user {
				result.Msg = "user not exists"
				result.Result = -2
				return
			}

			if password != user.PassToken {
				result.Msg = "invalid password"
				result.Result = -3
				return
			}

			//	done
			result.Result = 0
		}
	default:
		{
			result.Msg = "invalid ajax request"
		}
	}
}
