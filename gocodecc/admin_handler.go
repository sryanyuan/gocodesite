package gocodecc

import (
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
	}
}
