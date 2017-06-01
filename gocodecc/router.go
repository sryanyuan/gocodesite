package gocodecc

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cihub/seelog"
	"github.com/dchest/captcha"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

/*
	Permission
*/
const (
	kPermission_None       = iota // 默认权限，禁止访问
	kPermission_Guest             // 游客
	kPermission_User              // 注册用户
	kPermission_Admin             // 管理员
	kPermission_SuperAdmin        // 超级管理员
)

func checkPermission(perChecked uint32, want uint32) bool {
	if perChecked > kPermission_SuperAdmin ||
		want > kPermission_SuperAdmin {
		return false
	}

	if perChecked == kPermission_None ||
		want == kPermission_None {
		return false
	}

	if perChecked >= want {
		return true
	}

	return false
}

/*
	Http context
*/
type RequestContext struct {
	w         http.ResponseWriter
	r         *http.Request
	dbSession *sql.DB
	user      *WebUser
	tmRequest time.Time
	config    *AppConfig
}
type HttpHandler func(*RequestContext)

func (c *RequestContext) GetNginxRealIP() string {
	return c.r.Header.Get("X-real-ip")
}

func (c *RequestContext) GetRemoteIP() string {
	if c.config.NginxProxy {
		return c.GetNginxRealIP()
	}
	// Parse ip from remote addr
	remoteIPColonIndex := strings.LastIndex(c.r.RemoteAddr, ":")
	if -1 != remoteIPColonIndex {
		return c.r.RemoteAddr[:remoteIPColonIndex]
	}

	return ""
}

func (this *RequestContext) Redirect(url string, code int) {
	http.Redirect(this.w, this.r, url, code)
}

func (this *RequestContext) RenderJson(js interface{}) {
	renderJson(this, js)
}

func (this *RequestContext) RenderMessagePage(title string, text string, result bool) {
	renderMessage(this, title, text, result)
}

func (this *RequestContext) RenderDownloadPage(title string, text string, downloadUrl string) {
	url := fmt.Sprintf("/common/download?title=%s&text='%s'&url=%s", title, text, downloadUrl)
	this.Redirect(url, http.StatusFound)
}

func (this *RequestContext) RenderString(str string) {
	this.w.Write([]byte(str))
}

func (this *RequestContext) WriteHeader(header int) {
	this.w.WriteHeader(header)
}

func (this *RequestContext) WriteResponse(rsp []byte) (int, error) {
	return this.w.Write(rsp)
}

func (this *RequestContext) GetSession(name string) (*sessions.Session, error) {
	return store.Get(this.r, name)
}

func (this *RequestContext) GetWebUser() *WebUser {
	user := modelWebUserNew()
	session, err := this.GetSession("user")
	if nil != err {
		return user
	}

	userinfokey, ok := session.Values["login-key"].(string)
	if !ok {
		return user
	}

	//	parse info
	infoKeys := strings.Split(userinfokey, ":")
	if nil == infoKeys ||
		len(infoKeys) != 2 {
		return user
	}
	uid, err := strconv.Atoi(infoKeys[0])
	if nil != err ||
		0 == uid {
		return user
	}

	//	get user from db
	dbuser := modelWebUserGetUserByUid(uint32(uid))
	if dbuser.UserName != infoKeys[1] {
		return user
	}
	return dbuser
}

func (this *RequestContext) SaveWebUser(user *WebUser, saveDays int) {
	session, err := this.GetSession("user")
	if nil != err {
		return
	}

	if 0 == user.Uid {
		return
	}

	userinfokey := strconv.Itoa(int(user.Uid)) + ":" + user.UserName
	session.Values["login-key"] = userinfokey
	if 0 != saveDays {
		session.Options = &sessions.Options{
			MaxAge: saveDays * 24 * 60 * 60,
		}
	}
	session.Save(this.r, this.w)
}

func (this *RequestContext) ClearWebUser() {
	session, err := this.GetSession("user")
	if nil != err {
		return
	}

	session.Options = &sessions.Options{MaxAge: -1}
	session.Save(this.r, this.w)
}

/*
	Handler warper
*/
func responseWithAccessDenied(w http.ResponseWriter) {
	http.Error(w, "Access denied", http.StatusForbidden)
}

func wrapHandler(config *AppConfig, item *RouterItem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestCtx := RequestContext{
			w:         w,
			r:         r,
			dbSession: nil,
			tmRequest: time.Now(),
			config:    config,
		}

		user := requestCtx.GetWebUser()

		//	check permission
		if !checkPermission(user.Permission, item.Permission) {
			responseWithAccessDenied(w)
			return
		}

		seelog.Debug("Request url : ", r.URL)

		// Add site visitor counter
		var err error
		remoteIP := requestCtx.GetRemoteIP()
		if remoteIP == "" {
			seelog.Error("Get ip from request failed")
		} else {
			if err = modelSiteVisitorInc(remoteIP); nil != err {
				seelog.Error("Update site visitor failed:", err)
			}
		}

		requestCtx.user = user
		item.Handler(&requestCtx)
	}
}

/*
	Router item
*/
type RouterItem struct {
	Url        string      // 路由的url
	Permission uint32      // url访问权限
	Handler    HttpHandler // 处理器
}

var routerItems = []RouterItem{
	{"/", kPermission_Guest, indexHandler},
	{"/about", kPermission_Guest, aboutHander},
	{"/guestbook", kPermission_Guest, guestbookHandler},
	{"/account/signup", kPermission_Guest, signupHandler},
	{"/signin", kPermission_Guest, signinHandler},
	{"/signout", kPermission_User, signOutHandler},
	{"/articles", kPermission_Guest, articlesHandler},
	{"/mood", kPermission_Guest, moodHandler},
	{"/account/signupsuccess", kPermission_Guest, signupSuccessHandler},
	{"/member/{username}", kPermission_Guest, memberInfoHandler},
	{"/member/{username}/articles", kPermission_Guest, memberArticlesHandler},
	{"/project", kPermission_Guest, projectCategoryHandler},
	{"/project/{projectid:[0-9]*}/page/{page:[0-9]*}", kPermission_Guest, projectArticlesHandler},
	{"/project/{projectid:[0-9]*}/cmd/{cmd}", kPermission_Guest, projectArticleCmdHandler},
	{"/project/{projectid:[0-9]*}/article/{articleid:[0-9]*}", kPermission_Guest, projectArticleHandler},
	{"/ajax/{action}", kPermission_Guest, ajaxHandler},
	{"/admin/{action}", kPermission_SuperAdmin, adminHandler},
	{"/common/{action}", kPermission_Guest, commonHandler},
	{"/download/{filename}", kPermission_Guest, downloadHandler},
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[1:]
	http.ServeFile(w, r, filePath)
}

func InitRouters(config *AppConfig, r *mux.Router) {
	//	handle func
	routersCount := len(routerItems)
	for i := 0; i < routersCount; i++ {
		r.HandleFunc(routerItems[i].Url, wrapHandler(config, &routerItems[i]))
	}
	captchaStorage := captcha.NewMemoryStore(captcha.CollectNum, time.Minute*time.Duration(2))
	captcha.SetCustomStore(captchaStorage)
	http.Handle("/captcha/", captcha.Server(100, 40))

	//	static file
	http.Handle("/static/css/", http.FileServer(http.Dir(".")))
	http.Handle("/static/js/", http.FileServer(http.Dir(".")))
	http.Handle("/static/images/", http.FileServer(http.Dir(".")))
	http.Handle("/static/fonts/", http.FileServer(http.Dir(".")))
}
