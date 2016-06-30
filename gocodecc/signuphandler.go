package gocodecc

import (
	"net/http"
	"regexp"

	//"github.com/cihub/seelog"
)

type SignUpResult struct {
	Result int    `json:Result`
	Msg    string `json:Msg`
}

var signupRenderTpls = []string{
	"template/account/signup.tpl",
}

func signupHandler(ctx *RequestContext) {
	if ctx.user.Uid != 0 {
		//	already login
		ctx.Redirect("/", http.StatusOK)
		return
	}

	tplData := make(map[string]interface{})
	tplData["signup_result"] = ""
	tplData["signup_username"] = ""
	tplData["signup_email"] = ""

	//	render signup page
	if ctx.r.Method == "GET" {
		data := renderTemplate(ctx, signupRenderTpls, tplData)
		ctx.w.Write(data)
	} else if ctx.r.Method == "POST" {
		//	post register message
		ctx.r.ParseForm()
		failedMsg := ""

		//	validate input
		if matched, _ := regexp.Match("^[0-9a-zA-Z~!@$#%^]{5,20}$", []byte(ctx.r.Form.Get("user[password]"))); matched {
			if ctx.r.Form.Get("user[password_confirm]") != ctx.r.Form.Get("user[password]") {
				failedMsg = "两次输入密码不相同"
			}
		} else {
			failedMsg = "非法的密码"
		}
		if matched, _ := regexp.Match("^\\s*\\w+(?:\\.{0,1}[\\w-]+)*@[a-zA-Z0-9]+(?:[-.][a-zA-Z0-9]+)*\\.[a-zA-Z]+\\s*$", []byte(ctx.r.Form.Get("user[email]"))); matched {
			tplData["signup_email"] = ctx.r.Form.Get("user[email]")
		} else {
			failedMsg = "非法的邮件地址"
		}
		if matched, _ := regexp.Match("^([\u4E00-\uFA29]|[\uE7C7-\uE7F3]|[a-zA-Z0-9_]){4,10}$", []byte(ctx.r.Form.Get("user[name]"))); matched {
			tplData["signup_nickname"] = ctx.r.Form.Get("user[name]")
		} else {
			failedMsg = "非法的昵称"
		}
		if matched, _ := regexp.Match("^[a-zA-Z0-9_]{5,20}$", []byte(ctx.r.Form.Get("user[login]"))); matched {
			tplData["signup_username"] = ctx.r.Form.Get("user[login]")
		} else {
			failedMsg = "非法的用户名"
		}

		signUpResult := SignUpResult{
			Result: 1,
			Msg:    failedMsg,
		}
		if len(failedMsg) != 0 {
			//	echo error message
			renderJson(ctx, &signUpResult)
			return
		}

		//	all ok, redirect to signin page
		signUpResult.Result = 0
		signUpResult.Msg = "/"
		renderJson(ctx, &signUpResult)
	} else {
		http.Redirect(ctx.w, ctx.r, "/", http.StatusFound)
	}
}
