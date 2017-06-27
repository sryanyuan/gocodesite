package gocodecc

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

var (
	replyTableName = "reply"
	maxReplyLength = 512
)

type ReplyModel struct {
	Id         int `orm:pk;auto`
	Uid        uint32
	ReplyUser  string `orm:size(21)`
	Uri        string `orm:"size(256)"`
	IsDeleted  bool
	Comment    string `orm:"size(512)"`
	CreateTime int64
	UpdateTime int64
}

func (m *ReplyModel) TableName() string {
	return replyTableName
}

func init() {
	orm.RegisterModel(new(ReplyModel))
}

func modelReplyGetArticleReply(uri string, page int, limit int) ([]*ReplyModel, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}

	args := make([]interface{}, 0, 3)
	args = append(args, uri)
	sqlExpr := "SELECT id, uid, reply_user, is_deleted, comment, create_time, update_time FROM reply WHERE uri = ? ORDER BY create_time "
	if limit != 0 {
		sqlExpr += "LIMIT ? "
		args = append(args, limit)

		if page != 0 {
			sqlExpr += "OFFSET ? "
			args = append(args, page*limit)
		}
	}

	rows, err := db.Query(sqlExpr, args...)
	if nil != err {
		return nil, err
	}
	defer rows.Close()

	replys := make([]*ReplyModel, 0, 32)
	for rows.Next() {
		var reply ReplyModel

		if err = rows.Scan(&reply.Id,
			&reply.Uid,
			&reply.ReplyUser,
			&reply.IsDeleted,
			&reply.Comment,
			&reply.CreateTime,
			&reply.UpdateTime); nil != err {
			return nil, err
		}
		replys = append(replys, &reply)
	}

	return replys, nil
}

func modelReplyNew(uri string, user *WebUser, comment string) error {
	db, err := getRawDB()
	if nil != err {
		return err
	}

	if len(comment) > maxReplyLength {
		return errors.New("Reply max length is 512 characters")
	}
	if len(uri) == 0 {
		return errors.New("Invalid url for reply")
	}

	var reply ReplyModel
	reply.Comment = comment
	reply.CreateTime = time.Now().Unix()
	reply.IsDeleted = false
	reply.Uid = user.Uid
	reply.ReplyUser = user.UserName
	reply.Uri = uri

	if _, err = db.Exec("INSERT INTO reply (uid, reply_user, is_deleted, uri, comment, create_time) VALUES (?, ?, ?, ?, ?, ?)",
		reply.Uid, reply.ReplyUser, false, reply.Uri, reply.Comment, reply.CreateTime); nil != err {
		return err
	}
	return nil
}

func modelReplyGetCount(uri string) (int, error) {
	db, err := getRawDB()
	if nil != err {
		return 0, err
	}

	var cnt int
	row := db.QueryRow("SELECT COUNT(*) FROM reply WHERE uri = ?", uri)
	if err = row.Scan(&cnt); nil != err {
		return 0, err
	}

	return cnt, nil
}

func modelReplyMarkDelete(rid int) error {
	db, err := getRawDB()
	if nil != err {
		return err
	}

	_, err = db.Exec("UPDATE reply SET is_deleted = 1 WHERE id = ?", rid)
	return err
}

func modelReplyDelete(rid int) error {
	db, err := getRawDB()
	if nil != err {
		return err
	}

	_, err = db.Exec("DELETE FROM reply WHERE id = ?", rid)
	return err
}
