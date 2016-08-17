package gocodecc

var guestbookRenderTpls = []string{
	"template/guestbook.tpl",
	"template/component/comment_guestbook_html.tpl",
}

func guestbookHandler(ctx *RequestContext) {
	tplData := make(map[string]interface{})
	tplData["active"] = "guestbook"
	data := renderTemplate(ctx, guestbookRenderTpls, tplData)
	ctx.w.Write(data)
}
