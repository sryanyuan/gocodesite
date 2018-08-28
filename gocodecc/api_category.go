package gocodecc

import "net/http"

func init() {
	registerApi("/api/category", kPermission_Guest, apiCategoriesGet, []string{http.MethodGet})
}

type apiCategoryItem struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type apiCategoriesRsp struct {
	Categoreis []apiCategoryItem `json:"categories"`
}

func apiCategoriesGet(ctx *RequestContext) {
	cates, err := modelProjectCategoryGetAll()
	if nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	var rsp apiCategoriesRsp
	rsp.Categoreis = make([]apiCategoryItem, 0, len(cates))
	for _, cate := range cates {
		var item apiCategoryItem
		item.Name = cate.ProjectName
		item.Count = cate.ItemCount
		item.ID = cate.Id
		rsp.Categoreis = append(rsp.Categoreis, item)
	}
	ctx.WriteAPIRspOKWithMessage(&rsp)
}
