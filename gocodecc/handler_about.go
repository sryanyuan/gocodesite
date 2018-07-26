package gocodecc

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	aboutEditResumeTpls = []string{
		"template/about/edit_resume.html",
	}
)

func aboutHander(ctx *RequestContext) {
	dataCtx := map[string]interface{}{
		"active": "about",
	}
	dataHTML := renderTemplate(ctx, []string{"template/about.html"}, dataCtx)
	ctx.w.Write(dataHTML)
}

func aboutEditSectionHander(ctx *RequestContext) {
	vars := mux.Vars(ctx.r)
	section := vars["section"]
	switch section {
	case "resume":
		{
			if ctx.r.Method == http.MethodGet {
				// Load resume data
				var resumeData []byte
				var err error
				if "" != ctx.config.ResumeFile {
					resumeData, err = rawReadFileData(ctx.config.ResumeFile)
				}
				if nil != err {
					ctx.RenderMessagePage("错误", err.Error(), false)
					return
				}
				dataCtx := map[string]interface{}{
					"active":  "about",
					"content": string(resumeData),
				}
				data := renderTemplate(ctx, aboutEditResumeTpls, dataCtx)
				ctx.w.Write(data)
			} else if ctx.r.Method == http.MethodPost {
				aboutResumePost(ctx)
			}
		}
	default:
		{
			ctx.RenderMessagePage("错误", "Unknown section", false)
			return
		}
	}
}

func aboutResumePost(ctx *RequestContext) {
	var result AjaxResult
	result.Result = -1
	defer renderJson(ctx, &result)

	if ctx.user.Permission < kPermission_SuperAdmin {
		result.Msg = "Access denied"
		return
	}
	if ctx.config.ResumeFile == "" {
		result.Msg = "Resume file not set"
		return
	}
	ctx.r.ParseForm()
	contentMarkdown := ctx.r.Form.Get("editormd-markdown-doc")
	if len(contentMarkdown) == 0 {
		result.Msg = "Empty content"
		return
	}
	// Write into file
	readFileLock.Lock()
	err := ioutil.WriteFile(ctx.config.ResumeFile, []byte(contentMarkdown), 0644)
	if nil != err {
		delete(readFileCacheMap, ctx.config.ResumeFile)
	}
	readFileLock.Unlock()
	if nil != err {
		result.Msg = err.Error()
		return
	}
	result.Result = 0
	result.Msg = "/about"
}
