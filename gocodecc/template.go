package gocodecc

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"runtime"
	"strconv"
	"time"

	//"github.com/cihub/seelog"
)

var goVersion = runtime.Version()

var tplBinaryDataMap map[string]string

var tplFuncMap = template.FuncMap{
	"getprocesstime": tplfn_getprocesstime,
}

func init() {
	tplBinaryDataMap = make(map[string]string)
}

func tplfn_getprocesstime(tm time.Time) string {
	return strconv.Itoa((time.Now().Nanosecond()-tm.Nanosecond())/1e6) + " ms"
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
	data["imgPrefix"] = "/static/img"

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
