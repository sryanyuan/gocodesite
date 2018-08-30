package gocodecc

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
	//"github.com/cihub/seelog"
)

const (
	kUserNickNameLimit = 20
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
	ImportFrom    int
}

func (m *WebUser) TableName() string {
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

func modelWebUserGetSuperAdmin() (*WebUser, error) {
	var user WebUser
	o := orm.NewOrm()
	err := o.Raw("SELECT * FROM web_user WHERE permission = ?", kPermission_SuperAdmin).QueryRow(&user)
	if nil != err {
		return nil, err
	}
	return &user, nil
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

func modelWebUserGetCount() (int, error) {
	db, err := getRawDB()
	if nil != err {
		return 0, err
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM web_user").Scan(&count)
	if nil != err {
		return 0, err
	}
	return count, nil
}

func modelWebUserGetAll(limit, offset int) ([]*WebUser, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}

	users := make([]*WebUser, 0, 32)
	expr := "SELECT uid, permission, user_name, e_mail, last_login_time, create_time FROM web_user "
	if 0 != limit {
		expr += fmt.Sprintf("LIMIT %d ", limit)
	}
	if 0 != offset {
		expr += fmt.Sprintf("OFFSET %d", offset)
	}
	rows, err := db.Query(expr)
	if nil != err {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user WebUser
		if err = rows.Scan(&user.Uid,
			&user.Permission,
			&user.UserName,
			&user.EMail,
			&user.LastLoginTime,
			&user.CreateTime); nil != err {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}
