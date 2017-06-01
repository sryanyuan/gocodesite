package gocodecc

import (
	"database/sql"
	"strings"

	"github.com/astaxie/beego/orm"
	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var (
	dbDriverName = "sqlite3"
	dbDriverType = orm.DRSqlite
	dbDataSource = ""
)

func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// syncDB create all database and tables
func syncDB() error {
	var err error

	if err = orm.RunSyncdb("default", false, true); nil != err {
		return err
	}
	return nil
}

func initModels(ds string) error {
	var err error

	if err = orm.RegisterDriver(dbDriverName, dbDriverType); nil != err {
		return err
	}
	if err = orm.RegisterDataBase("default", dbDriverName, ds); nil != err {
		return err
	}

	return nil
}

func getRawDB() (*sql.DB, error) {
	db, err := orm.GetDB("default")
	return db, err
}
