package gocodecc

import (
	//"github.com/cihub/seelog"
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
			err := modelProjectCategoryUpdateProject(&originPrj)
			if nil != err {
				result.Msg = err.Error()
				return
			}

			//	done
			result.Result = 0
		}
	case "project_delete":
		{
			//	must be superadmin
			if ctx.user.Permission < kPermission_SuperAdmin {
				result.Msg = "permission denied"
				return
			}

			ctx.r.ParseForm()
			projectName := ctx.r.Form.Get("project[name]")

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
	default:
		{
			result.Msg = "invalid ajax request"
		}
	}
}
