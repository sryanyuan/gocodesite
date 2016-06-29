package gocodecc

import (
	"net/http"

	//"github.com/cihub/seelog"
)

func signupHandler(ctx *RequestContext) {
	if ctx.user.Uid != 0 {
		//	already login
		ctx.Redirect("/", http.StatusOK)
		return
	}

	tplData := make(map[string]interface{})
	tplData["signup_result"] = ""

	//	render signup page
	if ctx.r.Method == "GET" {
		data := renderTemplate(ctx, []string{"template/signup.tpl"}, tplData)
		ctx.w.Write(data)
	} else if ctx.r.Method == "POST" {
		//	post register message
		ctx.r.ParseForm()
		signUpSuccess := false

		//	if success, redirect to signup page
		if signUpSuccess {
			ctx.Redirect("/signup", http.StatusOK)
		} else {
			//	echo failed message
			failedMsg := "验证码错误"
			tplData["signup_result"] = failedMsg
			data := renderTemplate(ctx, []string{"template/signup.tpl"}, tplData)
			ctx.w.Write(data)
		}
	} else {
		http.Redirect(ctx.w, ctx.r, "/", http.StatusOK)
	}
}
