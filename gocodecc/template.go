package gocodecc

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"strconv"
)

var tplBinaryDataMap map[string]string

var tplFuncMap = template.FuncMap{
	"printhello": tplfn_printhello,
}

func init() {
	tplBinaryDataMap = make(map[string]string)
}

func tplfn_printhello(cnt int) string {
	return strconv.Itoa(cnt) + "hello"
}

func getTplBinaryData(file string) (string, error) {
	data, ok := tplBinaryDataMap[file]
	if ok {
		//	directy return data
		if !g_appConfig.Debug {
			return data
		}
	}

	layoutContent := ""
	layoutData, err := ioutil.ReadFile("/template/" + layoutFile)
	if nil != err {
		layoutContent = string(layoutData)
		tplBinaryDataMap[file] = layoutData
	}
	return layoutContent, err
}

func parseTemplate(fileName string, layoutFile string, data map[string]interface{}) []byte {
	var buffer bytes.Buffer
	t := template.New(layoutFile).Funcs(tplFuncMap)

	tplContent, err := getTplBinaryData(layoutFile)
	if nil != err {
		panic(err)
	}

	//	parse layout
	if err = t.Parse(tplContent); nil != err {
		panic(err)
	}

	//	parse file
	if t, err = t.ParseFiles("/template/" + fileName); nil != err {
		panic(err)
	}

	//	execute
	if err = t.Execute(&buffer, data); nil != err {
		panic(err)
	}

	return buffer.Bytes()
}
