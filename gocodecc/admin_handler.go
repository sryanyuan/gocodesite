package gocodecc

import (
	"path/filepath"

	"github.com/gorilla/mux"
)

var adminUploadRenderTpls []string = []string{
	"template/admin/upload.tpl",
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
	}
}
