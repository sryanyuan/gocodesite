package gocodecc

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"

	"github.com/cihub/seelog"
	"github.com/dchest/captcha"
)

var signinRenderTpls = []string{
	"template/account/signin.html",
}

type SignInResult struct {
	Result    int
	Msg       string
	CaptchaId string
}

func signinHandler(ctx *RequestContext) {
	tplData := make(map[string]interface{})
	if ctx.r.Method == "GET" {
		tplData["captchaid"] = captcha.NewLen(4)
		data := renderTemplate(ctx, signinRenderTpls, tplData)
		ctx.w.Write(data)
	} else {
		ctx.r.ParseForm()

		var result = SignInResult{
			Result: 1,
		}

		username := ctx.r.Form.Get("user[login]")
		password := ctx.r.Form.Get("user[password]")
		rememberMe := "0"
		url := ctx.r.Form.Get("url")
		if len(url) == 0 {
			url = "/"
		}

		//	check
		failedMsg := ""

		for {
			if !captcha.VerifyString(ctx.r.Form.Get("captchaid"), ctx.r.Form.Get("captchaSolution")) {
				failedMsg = "验证码错误"
				break
			}

			if len(username) == 0 {
				failedMsg = "用户名不能为空"
				break
			}
			if len(password) == 0 {
				failedMsg = "密码不能为空"
				break
			}
			if len(password) > 20 {
				failedMsg = "密码太长"
				break
			}

			// get user from db
			user := modelWebUserGetUserByUserName(username)
			if nil == user {
				failedMsg = "用户名不存在"
				break
			}
			md5calc := md5.New()
			md5calc.Write([]byte(password))
			md5Psw := hex.EncodeToString(md5calc.Sum(nil))
			if md5Psw != user.PassToken {
				failedMsg = "密码错误"
				break
			}

			//	now ok
			if "0" != rememberMe {
				ctx.SaveWebUser(user, 5)
			} else {
				ctx.SaveWebUser(user, 0)
			}
			break
		}

		if 0 != len(failedMsg) {
			result.CaptchaId = captcha.NewLen(4)
			result.Msg = failedMsg
			ctx.RenderJson(&result)
		} else {
			//	login ok
			result.Msg = url
			result.Result = 0
			ctx.RenderJson(&result)
			seelog.Debug("User ", username, " login success")
		}
	}
}

func signOutHandler(ctx *RequestContext) {
	//	already login
	if ctx.user.Uid == 0 {
		ctx.Redirect("/", http.StatusFound)
		return
	}

	//	clear user in session
	ctx.ClearWebUser()
	ctx.Redirect("/signin", http.StatusFound)
}
