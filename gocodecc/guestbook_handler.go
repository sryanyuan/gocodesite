package gocodecc

var guestbookRenderTpls = []string{
	"template/guestbook.tpl",
	"template/component/comment_embed.html",
	"template/component/comment_duoshuo.html",
	"template/component/comment_livere.html",
	"template/component/comment_163.html",
}

func guestbookHandler(ctx *RequestContext) {
	tplData := make(map[string]interface{})
	tplData["active"] = "guestbook"
	tplData["commentID"] = "guestbook"
	tplData["commentTitle"] = "留言板"
	data := renderTemplate(ctx, guestbookRenderTpls, tplData)
	ctx.w.Write(data)
}
