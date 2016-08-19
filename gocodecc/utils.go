package gocodecc

import (
	"bytes"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 获取文件大小的接口
type FileSize interface {
	Size() int64
}

// 获取文件信息的接口
type FileStat interface {
	Stat() (os.FileInfo, error)
}

//	trim html label
func trimHtmlLabel(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	return src
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//删除目录下的文件信息
//dirpath 目录路径
func delDirFile(dirpath string) error {
	//读取目录信息
	dir, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return err
	}

	for _, file := range dir {
		//读取到的是目录
		if file.IsDir() {
			continue
		}
		//文件的最后修改时间
		if err = os.Remove(dirpath + "/" + file.Name()); nil != err {
			return err
		}
	}

	return nil
}

func getOneImageFromHtml(html string) string {
	htmlReader := bytes.NewBuffer([]byte(html))
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if nil == err {
		imgNodes := doc.Find("img")
		if imgNodes.Length() != 0 {
			imgPath, exists := imgNodes.First().Attr("src")
			if exists {
				return imgPath
			}
			return ""
		}
	}
	return ""
}
