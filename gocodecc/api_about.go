package gocodecc

import (
	"io/ioutil"
	"net/http"
)

func init() {
	registerApi("/api/about", kPermission_Guest, apiAboutGet, []string{http.MethodGet})
	registerApi("/api/about/resume", kPermission_SuperAdmin, apiAboutResumePut, []string{http.MethodPut})
	registerApi("/api/resume/download", kPermission_SuperAdmin, apiResumeDownloadGet, []string{http.MethodGet})
}

type aboutGetContext struct {
	SiteIntro string `json:"siteintro"`
	Resume    string `json:"resume"`
}

func apiAboutGet(ctx *RequestContext) {
	mk := ctx.GetFormValueInt("mk", 0)
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
		if mk == 0 {
			gc.Resume, err = convertMarkdown2HTML(string(resumeData), 0)
			if nil != err {
				ctx.WriteAPIRspBadInternalError(err.Error())
				return
			}
		} else {
			gc.Resume = string(resumeData)
		}
	}
	ctx.WriteAPIRspOKWithMessage(&gc)
}

type apiAboutResumePutArg struct {
	Content string `json:"content"`
}

func apiAboutResumePut(ctx *RequestContext) {
	if ctx.config.ResumeFile == "" {
		ctx.WriteAPIRspBadInternalError("Resume file not set")
		return
	}
	var arg apiAboutResumePutArg
	if err := ctx.readFromBody(&arg); nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	if len(arg.Content) == 0 {
		ctx.WriteAPIRspBadInternalError("Content empty")
		return
	}
	// Write into file
	readFileLock.Lock()
	err := ioutil.WriteFile(ctx.config.ResumeFile, []byte(arg.Content), 0644)
	if nil == err {
		delete(readFileCacheMap, ctx.config.ResumeFile)
	}
	readFileLock.Unlock()
	if nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	ctx.WriteAPIRspOK(nil)
}

func apiResumeDownloadGet(ctx *RequestContext) {
	if ctx.config.ResumeFile == "" {
		ctx.WriteAPIRspBadInternalError("Resume file not set")
		return
	}
	resumeData, err := rawReadFileData(ctx.config.ResumeFile)
	if nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	ctx.w.Header().Set("Content-Type", "text/plain")
	ctx.w.Header().Set("Content-Disposition", "attachment;filename=resume.md")
	//ctx.w.Header().Set("Content-Length", len(fileBytes))
	ctx.w.Write(resumeData)
}
