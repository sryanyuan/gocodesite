package gocodecc

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type ArticleVistorModel struct {
	Id              int `orm:"pk;auto;index"`
	RecentVisitTime int64
	VisitTimes      int
	Uri             string `orm:"size(128)"`
	RemoteIp        string `orm:"size(21)"`
}

var (
	articleVistorTableName = "article_visitor"
)

func (m *ArticleVistorModel) TableName() string {
	return articleVistorTableName
}

func modelArticleVisitorNew(m *ArticleVistorModel) error {
	o := orm.NewOrm()
	m.RecentVisitTime = time.Now().Unix()
	_, err := o.Insert(m)
	return err
}

func modelArticleVisitorInc(uri string, ip string) error {
	db, err := getRawDB()
	if nil != err {
		return err
	}
	row := db.QueryRow("SELECT COUNT(*) FROM "+articleVistorTableName+" WHERE remote_ip = ? AND uri = ?", ip, uri)
	var cnt int
	if err = row.Scan(&cnt); nil != err {
		return err
	}
	if cnt == 0 {
		_, err = db.Exec("INSERT INTO "+articleVistorTableName+" (recent_visit_time, visit_times, remote_ip, uri) VALUES (?, ?, ?, ?)",
			time.Now().Unix(), 1, ip, uri)
		if nil != err {
			return err
		}
	} else {
		_, err = db.Exec("UPDATE "+articleVistorTableName+" SET visit_times = visit_times + 1 , recent_visit_time = ? WHERE remote_ip = ? AND uri = ?",
			time.Now().Unix(), ip, uri)
		if nil != err {
			return err
		}
	}
	return nil
}

func modelArticleVisitorGet(limit int) ([]*ArticleVistorModel, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}
	limitExp := ""
	if limit > 0 {
		limitExp = fmt.Sprintf(" LIMIT %d", limit)
	}
	rows, err := db.Query("SELECT recent_visit_time, visit_times, remote_ip, uri FROM " + articleVistorTableName + limitExp)
	if nil != err {
		return nil, err
	}
	defer rows.Close()

	results := make([]*ArticleVistorModel, 0, 128)
	for rows.Next() {
		var result ArticleVistorModel
		if err = rows.Scan(&result.RecentVisitTime, &result.VisitTimes, &result.RemoteIp, &result.Uri); nil != err {
			return nil, err
		}
		results = append(results, &result)
	}

	return results, nil
}

type SiteVisitorModel struct {
	Id              int `orm:"pk;auto;index"`
	RecentVisitTime int64
	VisitTimes      int
	RemoteIp        string `orm:"size(21);unique"`
}

var (
	siteVisitorTableName = "site_visitor"
)

func (m *SiteVisitorModel) TableName() string {
	return siteVisitorTableName
}

func modelSiteVisitorInc(ip string) error {
	db, err := getRawDB()
	if nil != err {
		return err
	}
	row := db.QueryRow("SELECT COUNT(*) FROM "+siteVisitorTableName+" WHERE remote_ip = ?", ip)
	var cnt int
	if err = row.Scan(&cnt); nil != err {
		return err
	}
	if cnt == 0 {
		_, err = db.Exec("INSERT INTO "+siteVisitorTableName+" (recent_visit_time, visit_times, remote_ip) VALUES (?, ?, ?)",
			time.Now().Unix(), 1, ip)
		if nil != err {
			return err
		}
	} else {
		_, err = db.Exec("UPDATE "+siteVisitorTableName+" SET visit_times = visit_times + 1 , recent_visit_time = ? WHERE remote_ip = ?",
			time.Now().Unix(), ip)
		if nil != err {
			return err
		}
	}
	return nil
}

func modelSiteVisitorGet(limit int) ([]*SiteVisitorModel, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}
	limitExp := ""
	if limit > 0 {
		limitExp = fmt.Sprintf(" LIMIT %d", limit)
	}
	rows, err := db.Query("SELECT recent_visit_time, visit_times, remote_ip FROM " + siteVisitorTableName + limitExp)
	if nil != err {
		return nil, err
	}
	defer rows.Close()

	results := make([]*SiteVisitorModel, 0, 128)
	for rows.Next() {
		var result SiteVisitorModel
		if err = rows.Scan(&result.RecentVisitTime, &result.VisitTimes, &result.RemoteIp); nil != err {
			return nil, err
		}
		results = append(results, &result)
	}

	return results, nil
}

func init() {
	orm.RegisterModel(new(ArticleVistorModel))
	orm.RegisterModel(new(SiteVisitorModel))
}
