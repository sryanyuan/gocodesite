package gocodecc

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"regexp"
	"time"

	"github.com/dchest/captcha"

	//"github.com/cihub/seelog"
)

type SignUpResult struct {
	Result    int    `json:Result`
	Msg       string `json:Msg`
	CaptchaId string `json:CaptchaId`
}

var signupRenderTpls = []string{
	"template/account/signup.tpl",
}

func signupHandler(ctx *RequestContext) {
	if ctx.user.Uid != 0 {
		//	already login
		ctx.Redirect("/", http.StatusFound)
		return
	}

	tplData := make(map[string]interface{})

	//	render signup page
	if ctx.r.Method == "GET" {
		tplData["captchaid"] = captcha.NewLen(4)
		data := renderTemplate(ctx, signupRenderTpls, tplData)
		ctx.w.Write(data)
	} else if ctx.r.Method == "POST" {
		//	post register message
		ctx.r.ParseForm()
		failedMsg := ""
		userName := ctx.r.Form.Get("user[login]")
		password := ctx.r.Form.Get("user[password]")
		email := ctx.r.Form.Get("user[email]")
		nickName := ctx.r.Form.Get("user[name]")

		//	validate input
		for {
			if !captcha.VerifyString(ctx.r.Form.Get("captchaid"), ctx.r.Form.Get("captchaSolution")) {
				failedMsg = "验证码错误"
				break
			}
			if matched, _ := regexp.Match("^[a-zA-Z0-9_]{5,20}$", []byte(userName)); !matched {
				failedMsg = "非法的用户名"
				break
			}
			if matched, _ := regexp.Match("^([\u4E00-\uFA29]|[\uE7C7-\uE7F3]|[a-zA-Z0-9_]){4,10}$", []byte(nickName)); !matched {
				failedMsg = "非法的昵称"
				break
			}
			if matched, _ := regexp.Match("^\\s*\\w+(?:\\.{0,1}[\\w-]+)*@[a-zA-Z0-9]+(?:[-.][a-zA-Z0-9]+)*\\.[a-zA-Z]+\\s*$", []byte(email)); !matched {
				failedMsg = "非法的邮件地址"
				break
			}
			if matched, _ := regexp.Match("^[0-9a-zA-Z~!@$#%^]{5,20}$", []byte(password)); matched {
				if ctx.r.Form.Get("user[password_confirm]") != password {
					failedMsg = "两次输入密码不相同"
					break
				}
			} else {
				failedMsg = "非法的密码"
				break
			}

			//	already exists?
			if userExists, _ := modelWebUserUserNameExists(userName); userExists {
				failedMsg = "用户名已存在"
				break
			}

			if nickExists, _ := modelWebUserNickNameExists(nickName); nickExists {
				failedMsg = "昵称已存在"
				break
			}
			break
		}

		signUpResult := SignUpResult{
			Result: 1,
			Msg:    failedMsg,
		}
		if len(failedMsg) != 0 {
			//	echo error message
			signUpResult.CaptchaId = captcha.NewLen(4)
			renderJson(ctx, &signUpResult)
			return
		}

		//	new user
		newuser := modelWebUserNew()
		newuser.CreateTime = time.Now().Unix()
		newuser.UserName = userName
		newuser.NickName = nickName

		md5calc := md5.New()
		md5calc.Write([]byte(password))
		newuser.PassToken = hex.EncodeToString(md5calc.Sum(nil))
		newuser.EMail = email

		if err := modelWebUserInsert(newuser); nil != err {
			signUpResult.CaptchaId = captcha.NewLen(4)
			signUpResult.Msg = kErrMsg_InternalError
			renderJson(ctx, &signUpResult)
			return
		}

		//	all ok, redirect to signin page
		signUpResult.Result = 0
		signUpResult.Msg = "/account/signupsuccess?account=" + userName
		renderJson(ctx, &signUpResult)
	} else {
		http.Redirect(ctx.w, ctx.r, "/", http.StatusFound)
	}
}

var signupSuccessRenderTpls = []string{
	"template/account/signupsuccess.tpl",
}

func signupSuccessHandler(ctx *RequestContext) {
	ctx.r.ParseForm()
	tplData := make(map[string]interface{})
	username := ctx.r.Form.Get("account")
	tplData["account"] = username
	data := renderTemplate(ctx, signupSuccessRenderTpls, tplData)
	ctx.w.Write(data)
}
