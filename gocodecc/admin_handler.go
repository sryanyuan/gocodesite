package gocodecc

import (
	"bytes"
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
			ctx.RenderDownloadPage("成功", "文件已打包入:"+zipPath, "/download/"+filepath.Base(zipPath)+"?t=markdown")
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
				tr := time.Unix(v.RecentVisitTime, 0)
				textBuffer.WriteString(tr.Format("2006-01-02 15:04:05"))
				textBuffer.WriteString("\r\n")
			}
			ctx.WriteResponse(textBuffer.Bytes())
		}
	case "site_visitors":
		{
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
				tr := time.Unix(v.RecentVisitTime, 0)
				textBuffer.WriteString(tr.Format("2006-01-02 15:04:05"))
				textBuffer.WriteString("\r\n")
			}
			ctx.WriteResponse(textBuffer.Bytes())
		}
	}
}
