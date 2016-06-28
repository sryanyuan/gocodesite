package gocodecc

func indexHandler(ctx *RequestContext) {
	dataCtx := map[string]interface{}{
		"active": "home",
	}
	dataHtml := renderTemplate(ctx, []string{"template/home.tpl"}, dataCtx)
	ctx.w.Write(dataHtml)
}
