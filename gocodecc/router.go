package gocodecc

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
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

/*
	Handler warper
*/
func responseWithAccessDenied(w http.ResponseWriter) {
	http.Error(w, "Access denied", http.StatusForbidden)
}

func getUserFromRequest(r *http.Request) *WebUser {
	//	get user
	var user WebUser

	//	not found, initialize as a guest
	user.Permission = kPermission_Guest
	user.Uid = 0
	user.UserName = "Guest"

	return &user
}

func wrapHandler(item *RouterItem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestCtx := RequestContext{
			w:         w,
			r:         r,
			dbSession: nil,
			tmRequest: time.Now(),
		}

		user := getUserFromRequest(r)

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

	//	static file
	http.Handle("/static/css/", http.FileServer(http.Dir(".")))
	http.Handle("/static/js/", http.FileServer(http.Dir(".")))
	http.Handle("/static/img/", http.FileServer(http.Dir(".")))
	http.Handle("/static/fonts/", http.FileServer(http.Dir(".")))
}
