package gocodecc

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cihub/seelog"
	"github.com/dchest/captcha"
	"github.com/gorilla/context"
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

// RequestContext wraps request and response
type RequestContext struct {
	w           http.ResponseWriter
	r           *http.Request
	dbSession   *sql.DB
	user        *WebUser
	tmRequest   time.Time
	config      *AppConfig
	requestBody []byte
}
type HttpHandler func(*RequestContext)

const (
	statusBad = 480
)

const (
	rspCodeInternalError = 1
)

type APIRsp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

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

func (c *RequestContext) Redirect(url string, code int) {
	http.Redirect(c.w, c.r, url, code)
}

func (c *RequestContext) RenderJson(js interface{}) {
	renderJson(c, js)
}

func (c *RequestContext) RenderMessagePage(title string, text string, result bool) {
	renderMessage(c, title, text, result)
}

func (c *RequestContext) RenderDownloadPage(title string, text string, downloadUrl string) {
	url := fmt.Sprintf("/common/download?title=%s&text='%s'&url=%s", title, text, downloadUrl)
	c.Redirect(url, http.StatusFound)
}

func (c *RequestContext) RenderString(str string) {
	c.w.Write([]byte(str))
}

func (c *RequestContext) WriteHeader(header int) {
	c.w.WriteHeader(header)
}

func (c *RequestContext) WriteResponse(rsp []byte) (int, error) {
	return c.w.Write(rsp)
}

func (c *RequestContext) GetSession(name string) (*sessions.Session, error) {
	return store.Get(c.r, name)
}

func (c *RequestContext) GetWebUser() *WebUser {
	user := modelWebUserNew()
	session, err := c.GetSession("user")
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

func (c *RequestContext) SaveWebUser(user *WebUser, saveDays int) {
	session, err := c.GetSession("user")
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
			Path:   "/",
			MaxAge: saveDays * 24 * 60 * 60,
		}
	}
	session.Save(c.r, c.w)
}

func (c *RequestContext) ClearWebUser() {
	session, err := c.GetSession("user")
	if nil != err {
		return
	}

	session.Options = &sessions.Options{
		Path:   "/",
		MaxAge: -1,
	}
	session.Save(c.r, c.w)
}

var defaultOKRsp APIRsp

func (c *RequestContext) WriteAPIRsp(header int, rsp *APIRsp) error {
	if header != http.StatusOK {
		seelog.Debugf("Rsp status bad with msg: %v", rsp)
	}
	if nil == rsp {
		rsp = &defaultOKRsp
	}

	jbytes, err := json.Marshal(rsp)
	if nil != err {
		return err
	}
	c.WriteHeader(header)
	c.WriteResponse(jbytes)
	return nil
}

func (c *RequestContext) WriteAPIRspOK(rsp *APIRsp) error {
	return c.WriteAPIRsp(http.StatusOK, rsp)
}

func (c *RequestContext) WriteAPIRspOKWithMessage(msg interface{}) error {
	var rsp APIRsp
	if nil != msg {
		jbytes, err := json.Marshal(msg)
		if nil != err {
			return err
		}
		rsp.Message = string(jbytes)
	}
	return c.WriteAPIRsp(http.StatusOK, &rsp)
}

func (c *RequestContext) WriteAPIRspBad(rsp *APIRsp) error {
	return c.WriteAPIRsp(statusBad, rsp)
}

func (c *RequestContext) WriteAPIRspBadInternalError(msg string) error {
	var rsp APIRsp
	rsp.Code = rspCodeInternalError
	rsp.Message = msg
	return c.WriteAPIRsp(statusBad, &rsp)
}

func (c *RequestContext) GetURLVarString(variable string) string {
	vars := mux.Vars(c.r)
	return strings.TrimSpace(vars[variable])
}

func (c *RequestContext) GetURLVarInt64(variable string, def int64) int64 {
	vars := mux.Vars(c.r)
	val := strings.TrimSpace(vars[variable])
	ival, err := strconv.ParseInt(val, 10, 64)
	if nil != err {
		return def
	}
	return ival
}

func (c *RequestContext) parseForm() error {
	return c.r.ParseForm()
}

func (c *RequestContext) GetFormValueString(key string) string {
	c.parseForm()
	return strings.Trim(c.r.Form.Get(key), " ")
}

func (c *RequestContext) GetFormValueInt(key string, def int) int {
	c.parseForm()
	val := c.GetFormValueString(key)
	if len(val) == 0 {
		// empty input
		return def
	}

	ival, err := strconv.Atoi(val)
	if nil != err {
		return def
	}

	// check max
	if ival > 0x7fffffff {
		return def
	}

	return ival
}

func (c *RequestContext) readBody() ([]byte, error) {
	if nil != c.requestBody {
		return c.requestBody, nil
	}

	c.r.ParseForm()
	data, err := ioutil.ReadAll(c.r.Body)
	if nil != err {
		return nil, err
	}
	c.requestBody = data

	return data, nil
}

func (c *RequestContext) readFromBody(i interface{}) error {
	body, err := c.readBody()
	if nil != err {
		return err
	}
	return json.Unmarshal(body, i)
}

/*
	Handler warper
*/
func responseWithAccessDenied(w http.ResponseWriter) {
	http.Error(w, "Access denied", http.StatusForbidden)
}

func wrapHandler(config *AppConfig, item *RouterItem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// NOTE: gorilla mux do not clear the request after go1.7
		defer context.Clear(r)

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
	Methods    []string
}

var routerItems = []RouterItem{
	{"/", kPermission_Guest, indexHandler, []string{http.MethodGet}},
	{"/about", kPermission_Guest, aboutHander, []string{http.MethodGet}},
	{"/about/edit/{section}", kPermission_SuperAdmin, aboutEditSectionHander, nil},
	{"/guestbook", kPermission_Guest, guestbookHandler, []string{http.MethodGet}},
	{"/donate", kPermission_Guest, donateHander, []string{http.MethodGet}},
	{"/donate/{orderid:[a-zA-Z0-9]*}", kPermission_Guest, donateCheckHandler, []string{http.MethodGet}},
	{"/account/signup", kPermission_Guest, signupHandler, nil},
	{"/signin", kPermission_Guest, signinHandler, nil},
	{"/signout", kPermission_User, signOutHandler, nil},
	{"/articles", kPermission_Guest, articlesHandler, []string{http.MethodGet}},
	{"/mood", kPermission_Guest, moodHandler, []string{http.MethodGet}},
	{"/account/signupsuccess", kPermission_Guest, signupSuccessHandler, []string{http.MethodGet}},
	{"/member/{username}", kPermission_Guest, memberInfoHandler, []string{http.MethodGet}},
	{"/member/{username}/articles", kPermission_Guest, memberArticlesHandler, []string{http.MethodGet}},
	{"/project", kPermission_Guest, projectCategoryHandler, []string{http.MethodGet}},
	{"/project/{projectid:[0-9]*}/page/{page:[0-9]*}", kPermission_Guest, projectArticlesHandler, []string{http.MethodGet}},
	{"/project/{projectid:[0-9]*}/cmd/{cmd}", kPermission_Guest, projectArticleCmdHandler, nil},
	{"/project/{projectid:[0-9]*}/article/{articleid:[0-9]*}/reply}", kPermission_User, projectArticleReplyHandler, []string{http.MethodGet}},
	{"/project/{projectid:[0-9]*}/article/{articleid:[0-9]*}", kPermission_Guest, projectArticleHandler, []string{http.MethodGet}},
	{"/ajax/{action}", kPermission_Guest, ajaxHandler, nil},
	{"/admin/{action}", kPermission_SuperAdmin, adminHandler, nil},
	{"/common/{action}", kPermission_Guest, commonHandler, nil},
	{"/download/{filename}", kPermission_Guest, downloadHandler, nil},
	{"/manager/{panel}", kPermission_SuperAdmin, managerPanelHandler, nil},
	{"/manager", kPermission_SuperAdmin, managerHandler, nil},
}

func registerRouter(path string, pem uint32, handler HttpHandler, methods []string) {
	routerItems = append(routerItems, RouterItem{
		Url:        path,
		Permission: pem,
		Handler:    handler,
		Methods:    methods,
	})
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[1:]
	http.ServeFile(w, r, filePath)
}

func InitRouters(config *AppConfig, r *mux.Router) {
	//	handle func
	routersCount := len(routerItems)
	for i := 0; i < routersCount; i++ {
		seelog.Debugf("Register router path %s, permission %d", routerItems[i].Url, routerItems[i].Permission)
		rt := r.HandleFunc(routerItems[i].Url, wrapHandler(config, &routerItems[i]))
		if nil != routerItems[i].Methods && 0 != len(routerItems[i].Methods) {
			rt.Methods(routerItems[i].Methods...)
		}
	}
	captchaStorage := captcha.NewMemoryStore(captcha.CollectNum, time.Minute*time.Duration(2))
	captcha.SetCustomStore(captchaStorage)
	captchaHandler := captcha.Server(100, 40)
	http.Handle("/captcha/", captchaHandler)
	http.Handle("/api/captcha/", captchaHandler)

	//	static file
	http.Handle("/static/css/", http.FileServer(http.Dir(".")))
	http.Handle("/static/js/", http.FileServer(http.Dir(".")))
	http.Handle("/static/images/", http.FileServer(http.Dir(".")))
	http.Handle("/static/fonts/", http.FileServer(http.Dir(".")))
	http.Handle("/static/img/", http.FileServer(http.Dir(".")))
	// New version frontend
	http.Handle("/view/", http.FileServer(http.Dir(".")))
}
