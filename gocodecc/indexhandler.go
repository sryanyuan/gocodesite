package gocodecc

var homeRenderTpls []string = []string{
	"template/component/article_detail_display.tpl",
	"template/home.tpl",
}

func indexHandler(ctx *RequestContext) {
	recentArticles, err := modelProjectArticleGetRecentArticles(8)
	if nil != err {
		ctx.RenderString("Internal error")
		return
	}

	dataCtx := map[string]interface{}{
		"active":         "home",
		"recentArticles": recentArticles,
	}
	dataHtml := renderTemplate(ctx, homeRenderTpls, dataCtx)
	ctx.w.Write(dataHtml)
}
