package gocodecc

import "github.com/dchest/captcha"

var guestbookRenderTpls = []string{
	"template/guestbook.html",
	"template/component/comment_embed.html",
	"template/component/comment_duoshuo.html",
	"template/component/comment_livere.html",
	"template/component/comment_163.html",
	"template/component/comment_native.html",
	"template/component/reply_list.html",
}

func guestbookHandler(ctx *RequestContext) {
	tplData := make(map[string]interface{})
	tplData["active"] = "guestbook"
	tplData["commentID"] = "guestbook"
	tplData["commentTitle"] = "留言板"
	var replies []*ReplyModel
	var err error
	if ctx.config.CommentProvider == "native" {
		tplData["captchaid"] = captcha.NewLen(4)
		// Replies
		replies, err = modelReplyGetArticleReply("/guestbook", 0, 0)
		if nil != err {
			ctx.RenderMessagePage("错误", err.Error(), false)
			return
		}
		tplData["replies"] = replies
	}

	data := renderTemplate(ctx, guestbookRenderTpls, tplData)
	ctx.w.Write(data)
}
