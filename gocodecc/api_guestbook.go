package gocodecc

import (
	"net/http"
)

const (
	guestbookCommentURI = "guestbook"
)

func init() {
	registerApi("/api/guestbook/comment", kPermission_Guest, apiGuestbookCommentsGet, []string{http.MethodGet})
	registerApi("/api/guestbook/comment", kPermission_User, apiArticleCommentPost, []string{http.MethodPost})
	registerApi("/api/guestbook/comment/{commentId}", kPermission_Guest, apiGuestbookCommentGet, []string{http.MethodGet})
}

func apiGuestbookCommentsGet(ctx *RequestContext) {
	rsp, err := getCommentsMergedByURI(guestbookCommentURI, ctx.user)
	if nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	ctx.WriteAPIRspOKWithMessage(rsp)
}

func apiGuestbookCommentGet(ctx *RequestContext) {
	commentId := int(ctx.GetURLVarInt64("commentId", 0))
	rsp, err := getCommentMergedByURI(guestbookCommentURI, commentId, ctx.user)
	if nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	ctx.WriteAPIRspOKWithMessage(rsp)
}
