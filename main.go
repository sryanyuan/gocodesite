package main

import (
	"flag"
	"log"
	"os"

	"os/exec"
	"path/filepath"

	"github.com/cihub/seelog"
	"github.com/sryanyuan/gocodesite/gocodecc"
)

func getModulePath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	dir := filepath.Dir(path)
	return dir
}

func initConf() bool {
	configpath := flag.String("configpath", "", "configpath <lua config file path>")
	flag.Parse()

	if len(*configpath) == 0 {
		flag.PrintDefaults()
		return false
	}

	return gocodecc.ReadLuaConfig(*configpath, "appconfig")
}

func main() {
	var err error
	//	initialize log
	logger, err := seelog.LoggerFromConfigAsFile(getModulePath() + "/static/conf/log.conf")
	if nil == logger {
		log.Println("Failed to initialize log module, error:", err)
		os.Exit(1)
	}
	seelog.ReplaceLogger(logger)

	//	clean up
	defer func() {
		e := recover()
		if nil != e {
			seelog.Error("Main routine quit with error:", e)
		} else {
			seelog.Info("Main routine quit normally")
		}

		seelog.Flush()
	}()

	//	load config
	if !initConf() {
		seelog.Error("Failed to init parameters")
		return
	}

	err = gocodecc.Start()
	if nil != err {
		seelog.Error("gocodecc error:", err)
	}
}
