package gocodecc

var guestbookRenderTpls = []string{
	"template/guestbook.tpl",
	"template/component/comment_guestbook_html.tpl",
	"template/component/comment_guestbook_html_duoshuo.tpl",
	"template/component/comment_guestbook_html_livere.tpl",
	"template/component/comment_guestbook_html_163.tpl",
}

func guestbookHandler(ctx *RequestContext) {
	tplData := make(map[string]interface{})
	tplData["active"] = "guestbook"
	data := renderTemplate(ctx, guestbookRenderTpls, tplData)
	ctx.w.Write(data)
}
