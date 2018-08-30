package gocodecc

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"regexp"
	"time"

	"github.com/dchest/captcha"
)

func init() {
	registerApi("/api/login/status", kPermission_Guest, apiLoginStatusGet, []string{http.MethodGet})
	registerApi("/api/login/captcha", kPermission_Guest, apiLoginCaptchaGet, []string{http.MethodGet})
	registerApi("/api/register", kPermission_Guest, apiRegisterPost, []string{http.MethodPost})
	registerApi("/api/login", kPermission_Guest, apiLoginPost, []string{http.MethodPost})
	registerApi("/api/logout", kPermission_Guest, apiLogoutPost, []string{http.MethodPost})
	registerApi("/api/superadmin", kPermission_Guest, apiSuperAdminGet, []string{http.MethodGet})
}

type loginStatusRsp struct {
	Role     int    `json:"role"`
	Uid      uint32 `json:"uid"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Sex      int    `json:"sex"`
	Mood     string `json:"mood"`
	Reason   string `json:"reason"`
}

func apiLoginStatusGet(ctx *RequestContext) {
	var rsp loginStatusRsp
	user := ctx.GetWebUser()
	if nil == user {
		rsp.Role = kPermission_Guest
	} else {
		rsp.Role = int(user.Permission)
		rsp.Username = user.UserName
		rsp.Uid = user.Uid
		rsp.Avatar = user.Avatar
		rsp.Sex = user.Sex
		rsp.Mood = user.Mood
	}
	ctx.WriteAPIRspOKWithMessage(&rsp)
}

type loginCaptchaRsp struct {
	Captcha string `json:"captcha"`
}

func apiLoginCaptchaGet(ctx *RequestContext) {
	var rsp loginCaptchaRsp
	rsp.Captcha = captcha.NewLen(4)
	ctx.WriteAPIRspOKWithMessage(&rsp)
}

type loginPostArg struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	CaptchaID  string `json:"captchaId"`
	Solution   string `json:"solution"`
	RememberMe bool   `json:"rememberMe"`
}

func apiLoginPost(ctx *RequestContext) {
	var arg loginPostArg
	if err := ctx.readFromBody(&arg); nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	if "" == arg.CaptchaID || "" == arg.Solution {
		ctx.WriteAPIRspBadInternalError("invalid captcha input")
		return
	}
	if !captcha.VerifyString(arg.CaptchaID, arg.Solution) {
		ctx.WriteAPIRspBadInternalError("invalid catpcha")
		return
	}
	if "" == arg.Username || "" == arg.Password {
		ctx.WriteAPIRspBadInternalError("invalid username or password")
		return
	}
	user := modelWebUserGetUserByUserName(arg.Username)
	if nil == user {
		ctx.WriteAPIRspBadInternalError("invalid username or password")
		return
	}
	md5calc := md5.New()
	md5calc.Write([]byte(arg.Password))
	md5Psw := hex.EncodeToString(md5calc.Sum(nil))
	if md5Psw != user.PassToken {
		ctx.WriteAPIRspBadInternalError("invalid username or password")
		return
	}

	// Remember me
	if arg.RememberMe {
		ctx.SaveWebUser(user, 5)
	} else {
		ctx.SaveWebUser(user, 0)
	}
	ctx.WriteAPIRspOK(nil)
}

func apiLogoutPost(ctx *RequestContext) {
	ctx.ClearWebUser()
	ctx.WriteAPIRspOK(nil)
}

type loginRegisterArg struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Mail      string `json:"mail"`
	CaptchaID string `json:"captchaId"`
	Solution  string `json:"solution"`
}

var (
	regAccount = regexp.MustCompile("^[a-zA-Z0-9_]{5,20}$")
	regMail    = regexp.MustCompile("^\\s*\\w+(?:\\.{0,1}[\\w-]+)*@[a-zA-Z0-9]+(?:[-.][a-zA-Z0-9]+)*\\.[a-zA-Z]+\\s*$")
	regPass    = regexp.MustCompile("^[0-9a-zA-Z~!@$#%^]{5,20}$")
)

func apiRegisterPost(ctx *RequestContext) {
	if ctx.user.Uid != 0 {
		ctx.WriteAPIRspBadInternalError("Already logined")
		return
	}
	var arg loginRegisterArg
	if err := ctx.readFromBody(&arg); nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	if "" == arg.CaptchaID || "" == arg.Solution {
		ctx.WriteAPIRspBadInternalError("invalid captcha input")
		return
	}
	if !captcha.VerifyString(arg.CaptchaID, arg.Solution) {
		ctx.WriteAPIRspBadInternalError("invalid catpcha")
		return
	}
	if !regAccount.MatchString(arg.Username) {
		ctx.WriteAPIRspBadInternalError("invalid username")
		return
	}
	if !regMail.MatchString(arg.Mail) {
		ctx.WriteAPIRspBadInternalError("invalid mail")
		return
	}
	if !regPass.MatchString(arg.Password) {
		ctx.WriteAPIRspBadInternalError("invalid password")
		return
	}

	newuser := modelWebUserNew()
	newuser.CreateTime = time.Now().Unix()
	newuser.UserName = arg.Username
	newuser.NickName = arg.Username
	newuser.Permission = kPermission_User

	md5calc := md5.New()
	md5calc.Write([]byte(arg.Password))
	newuser.PassToken = hex.EncodeToString(md5calc.Sum(nil))
	newuser.EMail = arg.Mail

	if err := modelWebUserInsert(newuser); nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}

	ctx.WriteAPIRspOK(nil)
}

func apiSuperAdminGet(ctx *RequestContext) {
	admin, err := modelWebUserGetSuperAdmin()
	if nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}

	var rsp loginStatusRsp
	rsp.Username = admin.UserName
	rsp.Uid = admin.Uid
	ctx.WriteAPIRspOKWithMessage(&rsp)
}
