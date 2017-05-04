package gocodecc

import "github.com/astaxie/beego/orm"
import "time"

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

func (this *ArticleVistorModel) TableName() string {
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
		_, err = db.Exec("UPDATE "+articleVistorTableName+" SET visit_times = visit_times + 1 WHERE remote_ip = ? AND uri = ?",
			ip, uri)
		if nil != err {
			return err
		}
	}
	return nil
}

type SiteVisitorModel struct {
	Id              int `orm:"pk;auto;index"`
	RecentVisitTime int64
	VisitTimes      int
	RemoteIp        string `orm:"size(21);unique"`
}

var (
	siteVistorTableName = "site_visitor"
)

func (this *SiteVisitorModel) TableName() string {
	return siteVistorTableName
}

func modelSiteVisitorInc(ip string) error {
	db, err := getRawDB()
	if nil != err {
		return err
	}
	row := db.QueryRow("SELECT COUNT(*) FROM "+siteVistorTableName+" WHERE remote_ip = ?", ip)
	var cnt int
	if err = row.Scan(&cnt); nil != err {
		return err
	}
	if cnt == 0 {
		_, err = db.Exec("INSERT INTO "+siteVistorTableName+" (recent_visit_time, visit_times, remote_ip) VALUES (?, ?, ?)",
			time.Now().Unix(), 1, ip)
		if nil != err {
			return err
		}
	} else {
		_, err = db.Exec("UPDATE "+siteVistorTableName+" SET visit_times = visit_times + 1 WHERE remote_ip = ?",
			ip)
		if nil != err {
			return err
		}
	}
	return nil
}

func init() {
	orm.RegisterModel(new(ArticleVistorModel))
	orm.RegisterModel(new(SiteVisitorModel))
}
