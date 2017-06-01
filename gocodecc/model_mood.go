package gocodecc

import (
	"database/sql"

	"github.com/astaxie/beego/orm"
)

type MoodInfo struct {
	Id       int `orm:"pk;auto;index"`
	PostTime int64
	Mood     string `orm:"size(1024)"`
	Image    string `orm:size(256)`
}

var (
	moodInfoTableName = "mood_info"
)

func init() {
	orm.RegisterModel(new(MoodInfo))
}

func (m *MoodInfo) TableName() string {
	return moodInfoTableName
}

func modelMoodInfoNew(info *MoodInfo) error {
	o := orm.NewOrm()
	_, err := o.Insert(info)
	return err
}

func modelMoodInfoDelete(id int) error {
	o := orm.NewOrm()
	var mood MoodInfo
	mood.Id = id
	_, err := o.Delete(&mood)
	return err
}

func modelMoodInfoUpdate(info *MoodInfo, cols []string) error {
	o := orm.NewOrm()
	_, err := o.Update(info, cols...)
	return err
}

func modelMoodInfoGet(page int, limit int) ([]*MoodInfo, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}

	var rows *sql.Rows
	if rows, err = db.Query(`SELECT 
	id,
	post_time,
	mood,
	image
	 FROM `+moodInfoTableName+" ORDER BY post_time DESC LIMIT ? OFFSET ?",
		limit,
		page*limit); nil != err {
		return nil, err
	}

	//	free the conn
	defer rows.Close()

	set := make([]*MoodInfo, 0, 10)

	for rows.Next() {
		var item MoodInfo
		err = rows.Scan(&item.Id, &item.PostTime, &item.Mood, &item.Image)
		if nil != err {
			return nil, err
		}
		set = append(set, &item)
	}

	return set, nil
}

func modelMoodGetCount() (int, error) {
	db, err := getRawDB()
	if nil != err {
		return 0, err
	}

	counter := 0
	err = db.QueryRow("SELECT COUNT(*) FROM " + moodInfoTableName).Scan(&counter)

	return counter, err
}
