package gocodecc

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var adminUploadRenderTpls []string = []string{
	"template/admin/upload.html",
}

func adminHandler(ctx *RequestContext) {
	vars := mux.Vars(ctx.r)
	action := vars["action"]

	switch action {
	case "upload":
		{
			data := renderTemplate(ctx, adminUploadRenderTpls, nil)
			ctx.w.Write(data)
		}
	case "resume_download":
		{
			if ctx.user.Permission < kPermission_SuperAdmin {
				ctx.RenderMessagePage("错误", "access denied", false)
				return
			}
			if ctx.config.ResumeFile == "" {
				ctx.RenderMessagePage("错误", "Resume file empty", false)
				return
			}
			f, err := os.Open(ctx.config.ResumeFile)
			if nil != err {
				ctx.RenderMessagePage("错误", err.Error(), false)
				return
			}
			defer f.Close()
			content, _ := ioutil.ReadAll(f)
			ctx.w.Header().Set("Content-Type", "text/plain")
			ctx.w.Header().Set("Content-Disposition", "attachment;filename="+path.Base(ctx.config.ResumeFile))
			ctx.w.Write(content)
		}
	case "pack_markdown":
		{
			if ctx.user.Permission < kPermission_SuperAdmin {
				ctx.RenderMessagePage("错误", "access denied", false)
				return
			}

			//	pack markdown
			zipPath, err := modelProjectArticlesPack("./markdown-articles/")
			if nil != err {
				ctx.RenderMessagePage("错误", err.Error(), false)
				return
			}
			ctx.RenderDownloadPage("成功", "文件已打包入:"+zipPath, "/download/"+filepath.Base(zipPath)+"?t=markdown_zip")
		}
	case "clean_markdown":
		{
			err := delDirFile("./markdown-articles/")
			if nil != err {
				ctx.RenderMessagePage("错误", err.Error(), false)
				return
			}
			ctx.RenderMessagePage("成功", "清理完毕", true)
		}
	case "article_visitors":
		{
			// Timezone
			tz, err := time.LoadLocation(siteTimezone)
			if nil != err {
				ctx.RenderMessagePage("错误", err.Error(), false)
				return
			}
			results, err := modelArticleVisitorGet(0)
			if nil != err {
				ctx.RenderMessagePage("错误", err.Error(), false)
				return
			}
			textBuffer := bytes.NewBuffer(nil)
			for _, v := range results {
				textBuffer.WriteString("IP:")
				textBuffer.WriteString(v.RemoteIp)
				textBuffer.WriteString(" \t| ")
				textBuffer.WriteString("URI:")
				textBuffer.WriteString(v.Uri)
				textBuffer.WriteString(" \t| ")
				textBuffer.WriteString("TIMES:")
				textBuffer.WriteString(strconv.Itoa(v.VisitTimes))
				textBuffer.WriteString(" \t| ")
				textBuffer.WriteString("RECENT:")
				tr := time.Unix(v.RecentVisitTime, 0).In(tz)
				textBuffer.WriteString(tr.Format("2006-01-02 15:04:05"))
				textBuffer.WriteString("\r\n")
			}
			ctx.WriteResponse(textBuffer.Bytes())
		}
	case "site_visitors":
		{
			// Timezone
			tz, err := time.LoadLocation(siteTimezone)
			if nil != err {
				ctx.RenderMessagePage("错误", err.Error(), false)
				return
			}
			results, err := modelSiteVisitorGet(0)
			if nil != err {
				ctx.RenderMessagePage("错误", err.Error(), false)
				return
			}
			textBuffer := bytes.NewBuffer(nil)
			for _, v := range results {
				textBuffer.WriteString("IP:")
				textBuffer.WriteString(v.RemoteIp)
				textBuffer.WriteString(" \t| ")
				textBuffer.WriteString("TIMES:")
				textBuffer.WriteString(strconv.Itoa(v.VisitTimes))
				textBuffer.WriteString(" \t| ")
				textBuffer.WriteString("RECENT:")
				tr := time.Unix(v.RecentVisitTime, 0).In(tz)
				textBuffer.WriteString(tr.Format("2006-01-02 15:04:05"))
				textBuffer.WriteString("\r\n")
			}
			ctx.WriteResponse(textBuffer.Bytes())
		}
	}
}
