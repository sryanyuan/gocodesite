package gocodecc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"sync"

	"github.com/cihub/seelog"
)

var goVersion = runtime.Version()

var tplBinaryDataMap map[string]string

var tplFuncMap = template.FuncMap{
	"getProcessTime":            tplfn_getprocesstime,
	"getUnixTimeString":         tplfn_getUnixTimeString,
	"getMemberAvatar":           tplfn_getMemberAvatar,
	"getTimeGapString":          tplfn_getTimeGapString,
	"articleEditable":           tplfn_articleEditable,
	"convertToHtml":             tplfn_convertToHtml,
	"getPageRange":              tplfn_getPageRange,
	"minusInt":                  tplfn_minusInt,
	"addInt":                    tplfn_addInt,
	"canPost":                   tplfn_canPost,
	"getThumb":                  tplfn_getThumb,
	"getImagePath":              tplfn_getImagePath,
	"formatDate":                tplfn_formatDate,
	"getArticleCoverImagePath":  tplfn_getArticleCoverImagePath,
	"getCategoryCoverImagePath": tplfn_getCategoryCoverImagePath,
	"getMoodImagePath":          tplfn_getMoodImagePath,
	"readFileData":              tplfn_readFileData,
}

func init() {
	tplBinaryDataMap = make(map[string]string)
}

func tplfn_getprocesstime(tm time.Time) string {
	return strconv.Itoa((time.Now().Nanosecond()-tm.Nanosecond())/1e6) + " ms"
}

func tplfn_getUnixTimeString(utm int64) string {
	tm := time.Unix(utm, 0)
	return tm.Format("2006-01-02")
}

func tplfn_getMemberAvatar(username string) string {
	user := modelWebUserGetUserByUserName(username)
	if nil == user {
		return "male.png"
	}

	if user.Sex == 0 {
		return "male.png"
	} else {
		return "female.png"
	}
}

func tplfn_articleEditable(user *WebUser, article *ProjectArticleItem) bool {
	if user.Permission > kPermission_Admin ||
		user.NickName == article.ArticleAuthor {
		if user.Uid != 0 {
			return true
		}
	}
	return false
}

func tplfn_getTimeGapString(tm int64) string {
	t := time.Unix(tm, 0)
	gap := time.Now().Sub(t)
	if gap.Seconds() < 60 {
		return "刚刚"
	}

	if gap.Minutes() < 60 {
		return fmt.Sprintf("%.0f 分钟前", gap.Minutes())
	}
	if gap.Hours() < 24 {
		return fmt.Sprintf("%.0f 小时前", gap.Hours())
	}

	hours := int(gap.Hours())
	days := hours / 24
	if days < 30 {
		return fmt.Sprintf("%d 天前", days)
	}
	// Get timezone
	tz, err := time.LoadLocation(siteTimezone)
	if nil != err {
		seelog.Error("Load localtion failed, timezone:", siteTimezone)
		return t.Format("2006-01-02 15:04")
	}
	return t.In(tz).Format("2006-01-02 15:04")
}

func tplfn_convertToHtml(str string) template.HTML {
	seelog.Info("Convert:", str)
	return template.HTML(str)
}

func tplfn_getPageRange(page int, showPage int, totalPage int) []int {
	pageStart := page - showPage/2
	if pageStart <= 0 {
		pageStart = 1
	}
	pageEnd := pageStart + showPage - 1

	pages := make([]int, 0, pageEnd-pageStart)
	for i := pageStart; i <= pageEnd; i++ {
		if i > totalPage {
			break
		}
		pages = append(pages, i)
	}

	pageCount := len(pages)
	if pageCount <= 0 {
		return pages
	}

	if pageCount < showPage &&
		pages[0] != 1 {
		//	need move
		offset := showPage - pageCount
		if pages[0]-offset >= 1 {
			prevPageBegin := pages[0] - offset
			prevPageEnd := pages[pageCount-1]
			pages = make([]int, 0, totalPage)
			for i := prevPageBegin; i <= prevPageEnd; i++ {
				pages = append(pages, i)
			}
		}
	}

	return pages
}

func tplfn_minusInt(val int, step int) int {
	return val - step
}

func tplfn_addInt(val int, step int) int {
	return val + step
}

func tplfn_getThumb(str string, charCount int) string {
	text := trimHtmlLabel(str)
	text = strings.TrimSpace(text)

	//	using rune
	runeText := []rune(text)
	if len(runeText) > charCount {
		text = string(runeText[:charCount])
		text += "..."
	}
	return text
}

func tplfn_canPost(cat *ProjectCategoryItem, user *WebUser) bool {
	if user.Uid == 0 {
		return false
	}

	if cat.Author == user.NickName {
		return true
	}

	if user.Permission >= cat.PostPriv {
		return true
	}
	return false
}

func tplfn_getImagePath(path string) string {
	path = strings.Trim(path, "/")
	path = strings.Trim(path, "\\")
	return kPrefixImagePath + "/" + path
}

func tplfn_formatDate(tm int64) string {
	timeVal := time.Unix(tm, 0)
	return timeVal.Format("2006-01-02")
}

func tplfn_getArticleCoverImagePath(projectId int, articleId int, path string) string {
	if len(path) == 0 {
		//	using default image
		return kPrefixImagePath + "/article_cover.png"
	}
	path = strings.Trim(path, "/")
	path = strings.Trim(path, "\\")
	return kPrefixImagePath + "/article-images/" + strconv.Itoa(projectId) + "/" + strconv.Itoa(articleId) + "/" + path
}

func tplfn_getCategoryCoverImagePath(path string) string {
	if len(path) == 0 {
		//	using default image
		return kPrefixImagePath + "/category_cover.png"
	}
	path = strings.Trim(path, "/")
	path = strings.Trim(path, "\\")
	return kPrefixImagePath + "/category-images/" + path
}

func tplfn_getMoodImagePath(path string) string {
	path = strings.Trim(path, "/")
	path = strings.Trim(path, "\\")
	return kPrefixImagePath + "/mood-images/" + path
}

var readFileCacheMap = make(map[string]template.HTML)
var readFileLock sync.Mutex

func tplfn_readFileData(path string) template.HTML {
	readFileLock.Lock()
	defer readFileLock.Unlock()

	// Find in cache
	if data, ok := readFileCacheMap[path]; ok {
		return data
	}

	// Open and cache
	f, err := os.Open(path)
	if nil != err {
		errMsg := fmt.Sprintf("Open file data error, file=%s, error=%s", path, err.Error())
		return template.HTML(errMsg)
	}
	defer f.Close()

	fileBytes, err := ioutil.ReadAll(f)
	if nil != err {
		errMsg := fmt.Sprintf("Read file data error, file=%s, error=%s", path, err.Error())
		seelog.Error(errMsg)
		return template.HTML(errMsg)
	}

	fileStr := template.HTML(fileBytes)
	readFileCacheMap[path] = fileStr

	return fileStr
}

func getTplBinaryData(file string, cache bool) string {
	data, ok := tplBinaryDataMap[file]
	if ok {
		//	directy return data
		if cache {
			return data
		}
	}

	layoutContent := ""
	layoutData, err := ioutil.ReadFile(file)
	if nil != err {
		panic(err)
	}

	layoutContent = string(layoutData)
	tplBinaryDataMap[file] = layoutContent

	return layoutContent
}

func parseTemplate(fileNames []string, cache bool, layoutFiles []string, data map[string]interface{}) []byte {
	var err error
	var buffer bytes.Buffer
	t := template.New("layout").Funcs(tplFuncMap)

	//	parse layout
	for _, v := range layoutFiles {
		tplContent := getTplBinaryData(v, cache)

		if t, err = t.Parse(tplContent); nil != err {
			panic(err)
		}
	}

	//	parse files
	if nil != fileNames &&
		len(fileNames) != 0 {
		if t, err = t.ParseFiles(fileNames...); nil != err {
			panic(err)
		}
	}

	//	execute
	if err = t.Execute(&buffer, data); nil != err {
		panic(err)
	}

	return buffer.Bytes()
}

func renderTemplate(rctx *RequestContext, fileNames []string, data map[string]interface{}) []byte {
	//	input some common variables
	if nil == data {
		data = make(map[string]interface{})
	}
	data["user"] = rctx.user

	_, ok := data["active"]
	if !ok {
		data["active"] = ""
	}

	data["goversion"] = goVersion
	data["requesttime"] = rctx.tmRequest
	data["config"] = rctx.config
	data["imgPrefix"] = kPrefixImagePath
	data["url"] = rctx.r.URL.Path

	//	get render data
	return parseTemplate(fileNames, !rctx.config.Debug, layoutFiles, data)
}

func renderJson(ctx *RequestContext, js interface{}) {
	if data, err := json.Marshal(js); nil != err {
		panic(err)
	} else {
		ctx.w.Header().Set("Content-Type", "application/json")
		ctx.w.Write(data)
	}
}

func renderMessage(ctx *RequestContext, title string, text string, ret bool) {
	result := ""
	if !ret {
		result = "1"
	}
	ctx.Redirect(fmt.Sprintf("/common/message?title=%s&text='%s'&result=%s", title, text, result), http.StatusFound)
}
