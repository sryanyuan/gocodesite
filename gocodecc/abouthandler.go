package gocodecc

func aboutHander(ctx *RequestContext) {
	dataCtx := map[string]interface{}{
		"active": "about",
	}
	dataHtml := renderTemplate(ctx, []string{"template/about.html"}, dataCtx)
	ctx.w.Write(dataHtml)
}
