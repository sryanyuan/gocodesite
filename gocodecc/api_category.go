package gocodecc

import (
	"net/http"
)

func init() {
	registerApi("/api/category", kPermission_Guest, apiCategoriesGet, []string{http.MethodGet})
	registerApi("/api/category", kPermission_SuperAdmin, apiCategoryPost, []string{http.MethodPost})
	registerApi("/api/category/{categoryId}", kPermission_SuperAdmin, apiCategoryGet, []string{http.MethodGet})
	registerApi("/api/category/{categoryId}", kPermission_SuperAdmin, apiCategoryPut, []string{http.MethodPut})
}

type apiCategoryRsp struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Count int    `json:"count"`
}

type apiCategoriesRsp struct {
	Categoreis []apiCategoryRsp `json:"categories"`
}

func apiCategoriesGet(ctx *RequestContext) {
	cates, err := modelProjectCategoryGetAll()
	if nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	var rsp apiCategoriesRsp
	rsp.Categoreis = make([]apiCategoryRsp, 0, len(cates))
	for _, cate := range cates {
		var item apiCategoryRsp
		item.Name = cate.ProjectName
		item.Count = cate.ItemCount
		item.ID = cate.Id
		item.Desc = cate.ProjectDescribe
		rsp.Categoreis = append(rsp.Categoreis, item)
	}
	ctx.WriteAPIRspOKWithMessage(&rsp)
}

type apiCategoryPostArg struct {
	Image string `json:"image"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Priv  int    `json:"priv"`
}

type apiCategoryPostRsp struct {
	CategoryId int    `json:"categoryId"`
	Name       string `json:"name"`
}

func apiCategoryPost(ctx *RequestContext) {
	var arg apiCategoryPostArg
	if err := ctx.readFromBody(&arg); nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	if len(arg.Name) == 0 ||
		len(arg.Desc) == 0 ||
		len(arg.Name) >= kCategoryNameLimit ||
		len(arg.Desc) >= kCategoryDescribeLimit {
		ctx.WriteAPIRspBadInternalError("invalid project name or project describe")
		return
	}

	var project ProjectCategoryItem
	project.Author = ctx.user.NickName
	project.Image = arg.Image
	project.ProjectName = arg.Name
	project.ProjectDescribe = arg.Desc
	project.PostPriv = uint32(kPermission_SuperAdmin)

	var categoryId int64
	var err error
	if categoryId, err = modelProjectCategoryAddReturnId(&project); nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}

	var rsp apiCategoryPostRsp
	rsp.CategoryId = int(categoryId)
	rsp.Name = arg.Name
	ctx.WriteAPIRspOKWithMessage(&rsp)
}

func apiCategoryGet(ctx *RequestContext) {
	categoryId := int(ctx.GetURLVarInt64("categoryId", 0))
	if 0 == categoryId {
		ctx.WriteAPIRspBadInternalError("invalid categoryId")
		return
	}
	var category ProjectCategoryItem
	err := modelProjectCategoryGetByProjectId(categoryId, &category)
	if nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	var rsp apiCategoryRsp
	rsp.Desc = category.ProjectDescribe
	rsp.Name = category.ProjectName
	ctx.WriteAPIRspOKWithMessage(&rsp)
}

func apiCategoryPut(ctx *RequestContext) {
	categoryId := int(ctx.GetURLVarInt64("categoryId", 0))
	if 0 == categoryId {
		ctx.WriteAPIRspBadInternalError("invalid categoryId")
		return
	}
	var arg apiCategoryPostArg
	if err := ctx.readFromBody(&arg); nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	var oldcate ProjectCategoryItem
	err := modelProjectCategoryGetByProjectId(categoryId, &oldcate)
	if nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	var newcate ProjectCategoryItem
	newcate = oldcate
	newcate.ProjectName = arg.Name
	newcate.ProjectDescribe = arg.Desc

	if len(newcate.ProjectName) == 0 ||
		len(newcate.ProjectDescribe) == 0 ||
		len(newcate.ProjectName) >= kCategoryNameLimit ||
		len(newcate.ProjectDescribe) >= kCategoryDescribeLimit {
		ctx.WriteAPIRspBadInternalError("invalid project name or project describe")
		return
	}
	if err = modelProjectCategoryUpdateProject(&oldcate, &newcate); nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}

	ctx.WriteAPIRspOK(nil)
}
