package gocodecc

import (
	"net/http"
	"time"

	"github.com/cihub/seelog"
)

func init() {
	registerApi("/api/archive", kPermission_Guest, apiArchiveGet, []string{http.MethodGet})
}

type apiArchiveGetRsp struct {
	Posts map[string][]*ProjectArticleItem `json:"posts"`
}

func apiArchiveGet(ctx *RequestContext) {
	articles, err := modelProjectArticleGetTitleAndPostTime()
	if nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	var rsp apiArchiveGetRsp
	rsp.Posts = make(map[string][]*ProjectArticleItem)
	for _, v := range articles {
		tm := time.Unix(v.PostTime, 0).Format(timeFormat)[0:7]
		plist, ok := rsp.Posts[tm]
		if !ok {
			plist = make([]*ProjectArticleItem, 0, 32)
			rsp.Posts[tm] = plist
			seelog.Info(tm)
		}
		plist = append(plist, v)
		rsp.Posts[tm] = plist
	}
	ctx.WriteAPIRspOKWithMessage(&rsp)
}
