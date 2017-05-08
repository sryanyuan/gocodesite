package gocodecc

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func downloadHandler(ctx *RequestContext) {
	vars := mux.Vars(ctx.r)
	filename := vars["filename"]
	ctx.r.ParseForm()
	fileType := ctx.r.Form.Get("t")
	fileType = strings.ToLower(fileType)

	switch fileType {
	case "markdown":
		{
			//	need super admin privilige
			if ctx.user.Permission < kPermission_SuperAdmin {
				ctx.RenderMessagePage("错误", "access denied", false)
				return
			}

			if len(filename) == 0 {
				ctx.RenderMessagePage("错误", "cannot find the file specific", false)
				return
			}

			//	open file
			f, err := os.Open("./markdown-articles/" + filename)
			if nil != err {
				ctx.RenderMessagePage("错误", "cannot open the file specific", false)
				return
			}
			defer f.Close()
			content, _ := ioutil.ReadAll(f)
			ctx.w.Header().Set("Content-Type", "application/zip")
			ctx.w.Write(content)
		}
	default:
		{
			ctx.RenderMessagePage("错误", "无效的文件索引符", false)
		}
	}
}
