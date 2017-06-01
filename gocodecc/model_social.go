package gocodecc

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/astaxie/beego/orm"
)

const (
	kSocialType_Weibo           = "Weibo"
	kSocialType_Github          = "Github"
	kSocialType_Facebook        = "Facebook"
	kSocialType_Twitter         = "Twitter"
	kSocialType_PersonalWebsite = "PersonalWebsite"
)

const (
	kSocialPrefix_Weibo           = "http://weibo.com"
	kSocialPrefix_Github          = "https://github.com"
	kSocialPrefix_Facebook        = "http://facebook.com"
	kSocialPrefix_Twitter         = "http://twitter.com"
	kSocialPrefix_PersonalWebsite = ""
)

type SocialInfo struct {
	Id              uint32 `orm:"pk;auto"`
	Uid             uint32 `orm:"index;unique"`
	Weibo           string `orm:"size(50)"`
	Github          string `orm:"size(50)"`
	Facebook        string `orm:"size(50)"`
	Twitter         string `orm:"size(50)"`
	PersonalWebsite string `orm:"size(50)"`
}

func (m *SocialInfo) TableName() string {
	return "social_info"
}

func init() {
	orm.RegisterModel(new(SocialInfo))
}

func rawSocialInfoExpandUrl(info *SocialInfo) {
	if len(info.Weibo) > 0 {
		info.Weibo = kSocialPrefix_Weibo + info.Weibo
	}
	if len(info.Github) > 0 {
		info.Github = kSocialPrefix_Github + info.Github
	}
	if len(info.Facebook) > 0 {
		info.Facebook = kSocialPrefix_Facebook + info.Facebook
	}
	if len(info.Twitter) > 0 {
		info.Twitter = kSocialPrefix_Twitter + info.Facebook
	}
}

func rawSocialInfoUpdateSocialWithType(info *SocialInfo, socialType string, value string) error {
	v := reflect.ValueOf(info)
	if !v.CanSet() ||
		v.Kind() != reflect.Ptr {
		return errors.New("Invalid info, not a pointer")
	}

	e := v.Elem()

	v = e.FieldByName(socialType)
	if v.IsNil() {
		return errors.New("Invalid field name:" + socialType)
	}
	v.Set(reflect.ValueOf(value))
	return nil
}

func modelSocialInfoExists(uid uint32) (uint32, error) {
	o := orm.NewOrm()

	var Id uint32
	err := o.Raw("SELECT Id FROM social_info WHERE Uid=?", uid).QueryRow(&Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return Id, nil
}

func modelSocialInfoGet(uid uint32) (*SocialInfo, error) {
	o := orm.NewOrm()
	var info SocialInfo

	err := o.Raw("SELECT * FROM social_info WHERE Uid=?", uid).QueryRow(&info)
	if err != nil {
		if err == sql.ErrNoRows {
			return &info, nil
		} else {
			return &info, err
		}
	}

	rawSocialInfoExpandUrl(&info)
	return &info, err
}

func modelSocialInfoUpdate(uid uint32, socialType string, value string) error {
	o := orm.NewOrm()

	if Id, err := modelSocialInfoExists(uid); nil != err {
		//	error
		return err
	} else {
		//	update or insert
		info := &SocialInfo{}

		if 0 == Id {
			//	insert
			rawSocialInfoUpdateSocialWithType(info, socialType, value)
			info.Id = Id
			info.Uid = uid
			if _, err = o.Insert(info); nil != err {
				return err
			}
		} else {
			//	update
			fieldName := snakeString(socialType)
			_, err := o.Raw("UPDATE social_info SET "+fieldName+" = ? WHERE Id = ?", value, Id).Exec()
			if nil != err {
				return err
			}
		}
	}

	return nil
}
