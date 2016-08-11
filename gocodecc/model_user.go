package gocodecc

import (
	"time"

	"github.com/astaxie/beego/orm"
	//"github.com/cihub/seelog"
)

type WebUser struct {
	Uid           uint32 `orm:"pk;auto;index"`
	Permission    uint32
	UserName      string `orm:"size(21);unique;index"`
	PassToken     string `orm:"size(32)"`
	Avatar        string `orm:"size(32)"`
	Sex           int
	NickName      string `orm:"size(20);unique"`
	EMail         string `orm:"size(50)"`
	AuthToken     string
	CreateTime    int64
	LastLoginTime int64
	MailVerified  bool
	LastLoginIp   string `orm:"size(15)"`
	Mood          string `orm:"size(128)"`
}

func (this *WebUser) TableName() string {
	return "web_user"
}

func init() {
	orm.RegisterModel(new(WebUser))
}

func modelWebUserNew() *WebUser {
	user := &WebUser{
		Uid:          0,
		Permission:   kPermission_Guest,
		UserName:     "Guest",
		Avatar:       "",
		CreateTime:   time.Now().Unix(),
		MailVerified: false,
	}
	return user
}

func modelWebUserInsert(user *WebUser) error {
	o := orm.NewOrm()
	_, err := o.Insert(user)
	return err
}

func modelWebUserUserNameExists(userName string) (bool, error) {
	o := orm.NewOrm()
	var user WebUser
	err := o.Raw("SELECT uid FROM web_user WHERE user_name = ?", userName).QueryRow(&user)
	if nil == err {
		return true, nil
	} else {
		if err == orm.ErrNoRows {
			return false, nil
		} else {
			return true, err
		}
	}
}

func modelWebUserNickNameExists(nickName string) (bool, error) {
	o := orm.NewOrm()
	var user WebUser
	err := o.Raw("SELECT uid FROM web_user WHERE nick_name=?", nickName).QueryRow(&user)
	if nil == err {
		return true, nil
	} else {
		if err == orm.ErrNoRows {
			return false, nil
		} else {
			return true, err
		}
	}
}

func modelWebUserGetUserByUid(uid uint32) *WebUser {
	user := modelWebUserNew()
	user.Uid = uid

	o := orm.NewOrm()
	if err := o.Read(user); nil != err {
		user.Uid = 0
	}
	return user
}

func modelWebUserGetUserByUserName(userName string) *WebUser {
	var user WebUser
	o := orm.NewOrm()
	err := o.Raw("SELECT * FROM web_user WHERE user_name=?", userName).QueryRow(&user)
	if nil != err {
		return nil
	}
	return &user
}
