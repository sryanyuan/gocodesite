package gocodecc

import (
	"database/sql"

	"github.com/astaxie/beego/orm"
)

var (
	bmkvTableName = "bmkv"
)

type BmkvModel struct {
	Id    int    `orm:pk;auto`
	Key   string `orm:"size(512);unique"`
	Value string `orm:"size(512)"`
}

func (m *BmkvModel) TableName() string {
	return bmkvTableName
}

func init() {
	orm.RegisterModel(new(BmkvModel))
}

func modelBmkvGet(key string) (string, error) {
	db, err := getRawDB()
	if nil != err {
		return "", err
	}

	row := db.QueryRow("SELECT value FROM bmkv WHERE key = ?", key)
	var value string
	if err = row.Scan(&value); nil != err {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

func modelBmkvGetAll(limit, offset int) ([]*BmkvModel, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}

	expr := "SELECT key, value FROM bmkv "
	args := make([]interface{}, 0, 2)

	if 0 != limit {
		expr += " LIMIT ? "
		args = append(args, limit)
	}
	if 0 != offset {
		expr += " OFFSET ? "
		args = append(args, offset)
	}

	rows, err := db.Query(expr, args...)
	if nil != err {
		return nil, err
	}
	defer rows.Close()

	res := make([]*BmkvModel, 0, 32)
	for rows.Next() {
		var kv BmkvModel
		if err = rows.Scan(&kv.Key, &kv.Value); nil != err {
			return nil, err
		}
		res = append(res, &kv)
	}
	return res, nil
}

func modelBmkvUpinsert(key, value string) error {
	db, err := getRawDB()
	if nil != err {
		return err
	}

	_, err = db.Exec("REPLACE INTO bmkv (key, value) VALUES(?, ?)", key, value)
	return err
}

func modelBmkvDelete(key string) error {
	db, err := getRawDB()
	if nil != err {
		return err
	}

	_, err = db.Exec("DELETE FROM bmkv WHERE key = ?", key)
	return err
}
