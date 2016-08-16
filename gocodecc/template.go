package gocodecc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var goVersion = runtime.Version()

var tplBinaryDataMap map[string]string

var tplFuncMap = template.FuncMap{
	"getProcessTime":    tplfn_getprocesstime,
	"getUnixTimeString": tplfn_getUnixTimeString,
	"getMemberAvatar":   tplfn_getMemberAvatar,
	"getTimeGapString":  tplfn_getTimeGapString,
	"articleEditable":   tplfn_articleEditable,
	"convertToHtml":     tplfn_convertToHtml,
	"getPageRange":      tplfn_getPageRange,
	"minusInt":          tplfn_minusInt,
	"addInt":            tplfn_addInt,
	"canPost":           tplfn_canPost,
	"getThumb":          tplfn_getThumb,
	"getImagePath":      tplfn_getImagePath,
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
	now := time.Now().Unix()
	gap := now - tm
	if gap < 0 {
		return "undefined"
	}

	year := gap / (365 * 30 * 24 * 60 * 60)
	if year > 0 {
		return fmt.Sprintf("%d 年前", year)
	}

	month := gap / (30 * 24 * 60 * 60)
	if month > 0 {
		return fmt.Sprintf("%d 月前", month)
	}

	day := gap / (24 * 60 * 60)
	if day > 0 {
		return fmt.Sprintf("%d 天前", day)
	}

	hour := gap / (60 * 60)
	if hour > 0 {
		return fmt.Sprintf("%d 小时前", hour)
	}

	minute := gap / 60
	if minute > 0 {
		return fmt.Sprintf("%d 分钟前", minute)
	}

	return fmt.Sprintf("%d 秒前", gap)
}

func tplfn_convertToHtml(str string) template.HTML {
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
	if len(text) > charCount {
		text = string(([]rune(text))[:charCount])
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

func getTplBinaryData(file string) string {
	data, ok := tplBinaryDataMap[file]
	if ok {
		//	directy return data
		if !g_appConfig.Debug {
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

func parseTemplate(fileNames []string, layoutFiles []string, data map[string]interface{}) []byte {
	var err error
	var buffer bytes.Buffer
	t := template.New("layout").Funcs(tplFuncMap)

	//	parse layout
	for _, v := range layoutFiles {
		tplContent := getTplBinaryData(v)

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
	data["config"] = &g_appConfig
	data["imgPrefix"] = kPrefixImagePath
	data["url"] = rctx.r.URL.Path

	//	get render data
	return parseTemplate(fileNames, layoutFiles, data)
}

func renderJson(ctx *RequestContext, js interface{}) {
	if data, err := json.Marshal(js); nil != err {
		panic(err)
	} else {
		ctx.w.Header().Set("Content-Type", "application/json")
		ctx.w.Write(data)
	}
}
