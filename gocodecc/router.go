package gocodecc

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	//"github.com/cihub/seelog"
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
}
type HttpHandler func(*RequestContext)

func (this *RequestContext) Redirect(url string, code int) {
	http.Redirect(this.w, this.r, url, code)
}

func (this *RequestContext) RenderJson(js interface{}) {
	renderJson(this, js)
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

func wrapHandler(item *RouterItem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestCtx := RequestContext{
			w:         w,
			r:         r,
			dbSession: nil,
			tmRequest: time.Now(),
		}

		user := requestCtx.GetWebUser()

		//	check permission
		if !checkPermission(user.Permission, item.Permission) {
			responseWithAccessDenied(w)
			return
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
	{"/account/signup", kPermission_Guest, signupHandler},
	{"/signin", kPermission_Guest, signinHandler},
	{"/signout", kPermission_User, signOutHandler},
	{"/account/signupsuccess", kPermission_Guest, signupSuccessHandler},
	{"/member/{username}", kPermission_Guest, memberInfoHandler},
	{"/project", kPermission_Guest, projectCategoryHandler},
	{"/project/{projectname}/page/{page:[0-9]*}", kPermission_Guest, projectArticlesHandler},
	{"/project/{projectname}/cmd/{cmd}", kPermission_SuperAdmin, projectArticleCmdHandler},
	{"/project/{projectname}/article/{articleid:[0-9]*}", kPermission_Guest, projectArticleHandler},
	{"/ajax/{action}", kPermission_Guest, ajaxHandler},
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[1:]
	http.ServeFile(w, r, filePath)
}

func InitRouters(r *mux.Router) {
	//	handle func
	routersCount := len(routerItems)
	for i := 0; i < routersCount; i++ {
		r.HandleFunc(routerItems[i].Url, wrapHandler(&routerItems[i]))
	}
	captchaStorage := captcha.NewMemoryStore(captcha.CollectNum, time.Minute*time.Duration(2))
	captcha.SetCustomStore(captchaStorage)
	http.Handle("/captcha/", captcha.Server(100, 40))

	//	static file
	http.Handle("/static/css/", http.FileServer(http.Dir(".")))
	http.Handle("/static/js/", http.FileServer(http.Dir(".")))
	//http.Handle("/static/img/", http.FileServer(http.Dir(".")))
	http.Handle("/static/images/", http.FileServer(http.Dir(".")))
	http.Handle("/static/fonts/", http.FileServer(http.Dir(".")))
}
