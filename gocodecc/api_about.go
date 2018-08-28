package gocodecc

import (
	"net/http"
)

func init() {
	registerApi("/api/about", kPermission_Guest, apiAboutGet, []string{http.MethodGet})
	registerApi("/api/about/resume", kPermission_SuperAdmin, apiAboutGet, []string{http.MethodPost})
}

type aboutGetContext struct {
	SiteIntro string `json:"siteintro"`
	Resume    string `json:"resume"`
}

func apiAboutGet(ctx *RequestContext) {
	var gc aboutGetContext
	if ctx.config.AboutHTMLFile != "" {
		htmlData, err := rawReadFileData(ctx.config.AboutHTMLFile)
		if nil != err {
			ctx.WriteAPIRspBad(&APIRsp{
				Code:    1,
				Message: err.Error(),
			})
			return
		}
		gc.SiteIntro = string(htmlData)
	}
	if ctx.config.ResumeFile != "" {
		resumeData, err := rawReadFileData(ctx.config.ResumeFile)
		if nil != err {
			ctx.WriteAPIRspBad(&APIRsp{
				Code:    1,
				Message: err.Error(),
			})
			return
		}
		gc.Resume, err = convertMarkdown2HTML(string(resumeData), 0)
		if nil != err {
			ctx.WriteAPIRspBadInternalError(err.Error())
			return
		}
	}
	ctx.WriteAPIRspOKWithMessage(&gc)
}

func apiAboutResumePost(ctx *RequestContext) {

}
