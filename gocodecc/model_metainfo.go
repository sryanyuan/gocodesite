package gocodecc

import (
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

type MetaInfo struct {
	Id    int    `orm:"pk;auto;index"`
	Key   string `orm:"size(128);index;unique"`
	Value string `orm:"size(128)"`
}

var (
	metaInfoTableName = "meta_info"
)

func init() {
	orm.RegisterModel(new(MetaInfo))
}

func (this *MetaInfo) TableName() string {
	return metaInfoTableName
}

func initMetaInfo() {
	o := orm.NewOrm()

	createSiteMetaInfo := &MetaInfo{
		Key:   "create_site_time",
		Value: strconv.FormatInt(time.Now().Unix(), 10),
	}
	o.Insert(createSiteMetaInfo)
}

func modelMetaInfoGet(key string) (string, error) {
	o := orm.NewOrm()
	metaInfo := MetaInfo{
		Key: key,
	}

	err := o.Read(&metaInfo, "Key")
	if nil != err {
		return "", err
	}
	return metaInfo.Value, nil
}
