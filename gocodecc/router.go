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
	ri          *RouterItem
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
	rspCodeNeedLogin
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

func (c *RequestContext) WriteAPIRspBadNeedLogin(msg string) error {
	var rsp APIRsp
	rsp.Code = rspCodeNeedLogin
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
			ri:        item,
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
// RouterMeta should be initialized before handler invoke
type RouterMeta struct {
	vals map[string]interface{}
}

func (m *RouterMeta) GetInt(key string) (int, bool) {
	if nil == m.vals {
		return 0, false
	}
	v, ok := m.vals[key]
	if !ok {
		return 0, false
	}
	iv, ok := v.(int)
	if !ok {
		return 0, false
	}
	return iv, true
}

func (m *RouterMeta) Init(mk map[string]interface{}) {
	m.vals = mk
}

type RouterItem struct {
	Url        string      // 路由的url
	Permission uint32      // url访问权限
	Handler    HttpHandler // 处理器
	Methods    []string
	Meta       RouterMeta
}

var routerItems = []RouterItem{
	{
		Url:        "/",
		Permission: kPermission_Guest,
		Handler:    indexHandler,
		Methods:    []string{http.MethodGet},
	},
	{
		Url:        "/about",
		Permission: kPermission_Guest,
		Handler:    aboutHander,
		Methods:    []string{http.MethodGet}},
	{
		Url:        "/about/edit/{section}",
		Permission: kPermission_SuperAdmin,
		Handler:    aboutEditSectionHander,
	},
	{
		Url:        "/guestbook",
		Permission: kPermission_Guest,
		Handler:    guestbookHandler,
		Methods:    []string{http.MethodGet},
	},
	{
		Url:        "/donate",
		Permission: kPermission_Guest,
		Handler:    donateHander,
		Methods:    []string{http.MethodGet},
	},
	{
		Url:        "/donate/{orderid:[a-zA-Z0-9]*}",
		Permission: kPermission_Guest,
		Handler:    donateCheckHandler,
		Methods:    []string{http.MethodGet},
	},
	{
		Url:        "/account/signup",
		Permission: kPermission_Guest,
		Handler:    signupHandler,
	},
	{
		Url:        "/signin",
		Permission: kPermission_Guest,
		Handler:    signinHandler,
	},
	{
		Url:        "/signout",
		Permission: kPermission_User,
		Handler:    signOutHandler,
	},
	{
		Url:        "/articles",
		Permission: kPermission_Guest,
		Handler:    articlesHandler,
		Methods:    []string{http.MethodGet},
	},
	{
		Url:        "/mood",
		Permission: kPermission_Guest,
		Handler:    moodHandler,
		Methods:    []string{http.MethodGet},
	},
	{
		Url:        "/account/signupsuccess",
		Permission: kPermission_Guest,
		Handler:    signupSuccessHandler,
		Methods:    []string{http.MethodGet},
	},
	{
		Url:        "/member/{username}",
		Permission: kPermission_Guest,
		Handler:    memberInfoHandler,
		Methods:    []string{http.MethodGet},
	},
	{
		Url:        "/member/{username}/articles",
		Permission: kPermission_Guest,
		Handler:    memberArticlesHandler,
		Methods:    []string{http.MethodGet},
	},
	{
		Url:        "/project",
		Permission: kPermission_Guest,
		Handler:    projectCategoryHandler,
		Methods:    []string{http.MethodGet},
	},
	{
		Url:        "/project/{projectid:[0-9]*}/page/{page:[0-9]*}",
		Permission: kPermission_Guest,
		Handler:    projectArticlesHandler,
		Methods:    []string{http.MethodGet},
	},
	{
		Url:        "/project/{projectid:[0-9]*}/cmd/{cmd}",
		Permission: kPermission_Guest,
		Handler:    projectArticleCmdHandler,
	},
	{
		Url:        "/project/{projectid:[0-9]*}/article/{articleid:[0-9]*}/reply}",
		Permission: kPermission_User,
		Handler:    projectArticleReplyHandler,
		Methods:    []string{http.MethodGet},
	},
	{
		Url:        "/project/{projectid:[0-9]*}/article/{articleid:[0-9]*}",
		Permission: kPermission_Guest,
		Handler:    projectArticleHandler,
		Methods:    []string{http.MethodGet},
	},
	{
		Url:        "/ajax/{action}",
		Permission: kPermission_Guest,
		Handler:    ajaxHandler,
	},
	{
		Url:        "/admin/{action}",
		Permission: kPermission_SuperAdmin,
		Handler:    adminHandler,
	},
	{
		Url:        "/common/{action}",
		Permission: kPermission_Guest,
		Handler:    commonHandler,
	},
	{
		Url:        "/download/{filename}",
		Permission: kPermission_Guest,
		Handler:    downloadHandler,
	},
	{
		Url:        "/manager/{panel}",
		Permission: kPermission_SuperAdmin,
		Handler:    managerPanelHandler,
	},
	{
		Url:        "/manager",
		Permission: kPermission_SuperAdmin,
		Handler:    managerHandler,
	},
}

func registerRouter(path string, pem uint32, handler HttpHandler, methods []string, meta map[string]interface{}) {
	ri := RouterItem{
		Url:        path,
		Permission: pem,
		Handler:    handler,
		Methods:    methods,
	}
	if nil != meta {
		ri.Meta.Init(meta)
	}
	routerItems = append(routerItems, ri)
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
