package gocodecc

var homeRenderTpls []string = []string{
	"template/component/article_detail_display.tpl",
	"template/home.tpl",
}

func indexHandler(ctx *RequestContext) {
	//	get recent articles
	recentArticles, err := modelProjectArticleGetRecentArticles(0, 5)
	if nil != err {
		ctx.RenderMessagePage("错误", err.Error(), false)
		return
	}

	//	get article count
	articleCount, err := modelProjectArticleGetArticleCountAll()
	if nil != err {
		ctx.RenderMessagePage("错误", err.Error(), false)
		return
	}

	//	get category
	category, err := modelProjectCategoryGetAllSimple()
	if nil != err {
		ctx.RenderMessagePage("错误", err.Error(), false)
		return
	}

	//	get member count
	memberCount, err := modelWebUserGetCount()
	if nil != err {
		ctx.RenderMessagePage("错误", err.Error(), false)
		return
	}

	dataCtx := map[string]interface{}{
		"active":         "home",
		"recentArticles": recentArticles,
		"articleCount":   articleCount,
		"category":       category,
		"memberCount":    memberCount,
		"createSiteTime": metaInfoCreateSiteTime,
	}
	dataHtml := renderTemplate(ctx, homeRenderTpls, dataCtx)
	ctx.w.Write(dataHtml)
}
