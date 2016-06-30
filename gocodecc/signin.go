package gocodecc

var signinRenderTpls = []string{
	"template/account/signin.tpl",
}

type SignInResult struct {
	Result int
	Msg    string
}

func signinHandler(ctx *RequestContext) {
	tplData := make(map[string]interface{})
	if ctx.r.Method == "GET" {
		data := renderTemplate(ctx, signinRenderTpls, tplData)
		ctx.w.Write(data)
	} else {
		ctx.r.ParseForm()

		var result = SignInResult{
			Result: 1,
		}

		username := ctx.r.Form.Get("user[login]")
		password := ctx.r.Form.Get("user[password]")

		//	check
		failedMsg := ""

		if len(username) == 0 {
			failedMsg = "用户名不能为空"
		}
		if len(password) == 0 {
			failedMsg = "密码不能为空"
		}
		if len(password) > 20 {
			failedMsg = "密码太长"
		}

		if 0 != len(failedMsg) {
			result.Msg = failedMsg
			ctx.RenderJson(&result)
		} else {
			result.Msg = "/"
			result.Result = 0
			ctx.RenderJson(&result)
		}
	}
}
