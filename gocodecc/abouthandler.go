package gocodecc

func aboutHander(ctx *RequestContext) {
	dataCtx := map[string]interface{}{
		"active": "about",
	}
	dataHTML := renderTemplate(ctx, []string{"template/about.html"}, dataCtx)
	ctx.w.Write(dataHTML)
}
