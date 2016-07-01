package gocodecc

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func initModels() error {
	var err error

	if err = orm.RegisterDriver("mysql", orm.DRMySQL); nil != err {
		return err
	}
	if err = orm.RegisterDataBase("default", "mysql", g_appConfig.DBAddress); nil != err {
		return err
	}
	if err = orm.RunSyncdb("default", false, g_appConfig.Debug); nil != err {
		return err
	}

	return nil
}
