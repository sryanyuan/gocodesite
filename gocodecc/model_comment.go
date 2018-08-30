package gocodecc

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
)

var (
	commentTableName = "comment"
	maxCommentLength = 512
)

type CommentModel struct {
	Id         int `orm:pk;auto`
	Uid        uint32
	ReplyUser  string `orm:size(21)`
	Uri        string `orm:"size(256)"`
	IsSub      bool
	SubRefId   int
	SubToUid   uint32
	SubToUser  string `orm:"size(21)"`
	Comment    string `orm:"size(512)"`
	CreateTime int64
	UpdateTime int64
	Agree      int
	Oppose     int
	Review     int
}

func (m *CommentModel) TableName() string {
	return commentTableName
}

func init() {
	orm.RegisterModel(new(CommentModel))
}

func modelCommentGet(id int) (*CommentModel, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}

	sqlExpr := `SELECT 
				uid,
				reply_user, 
				uri,
				is_sub, 
				sub_ref_id, 
				sub_to_uid,
				sub_to_user, 
				comment, 
				create_time, 
				update_time,
				agree,
				oppose,
				review FROM comment WHERE id = ?`

	row := db.QueryRow(sqlExpr, id)
	var reply CommentModel

	if err = row.Scan(
		&reply.Uid,
		&reply.ReplyUser,
		&reply.Uri,
		&reply.IsSub,
		&reply.SubRefId,
		&reply.SubToUid,
		&reply.SubToUser,
		&reply.Comment,
		&reply.CreateTime,
		&reply.UpdateTime,
		&reply.Agree,
		&reply.Oppose,
		&reply.Review); nil != err {
		return nil, err
	}
	reply.Id = id

	return &reply, nil
}

func modelCommentGetTopCount(uri string, all bool) (int, error) {
	db, err := getRawDB()
	if nil != err {
		return 0, err
	}
	expr := "SELECT COUNT(*) FROM " + commentTableName + " WHERE uri = ?"
	if !all {
		expr += " AND review = 1 "
	}
	row := db.QueryRow(expr, uri)
	var cnt int
	if err := row.Scan(&cnt); nil != err {
		return 0, err
	}
	return cnt, nil
}

func modelCommentGetSubs(uri string, subRefID int, page, limit int, all bool) ([]*CommentModel, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}

	args := make([]interface{}, 0, 3)
	args = append(args, uri, subRefID)
	sqlExpr := `SELECT 
				id, 
				uid,
				reply_user, 
				is_sub, 
				sub_ref_id, 
				sub_to_uid,
				sub_to_user, 
				comment, 
				create_time, 
				update_time,
				agree,
				oppose,
				review FROM comment WHERE uri = ? AND sub_ref_id = ?`
	if !all {
		sqlExpr += " AND review = 1 "
	}
	sqlExpr += " ORDER BY create_time "
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

	replys := make([]*CommentModel, 0, 32)
	for rows.Next() {
		var reply CommentModel

		if err = rows.Scan(&reply.Id,
			&reply.Uid,
			&reply.ReplyUser,
			&reply.IsSub,
			&reply.SubRefId,
			&reply.SubToUid,
			&reply.SubToUser,
			&reply.Comment,
			&reply.CreateTime,
			&reply.UpdateTime,
			&reply.Agree,
			&reply.Oppose,
			&reply.Review); nil != err {
			return nil, err
		}
		replys = append(replys, &reply)
	}

	return replys, nil
}

func modelCommentGetAllUnreviewedCount() (int, error) {
	db, err := getRawDB()
	if nil != err {
		return 0, err
	}

	row := db.QueryRow("SELECT COUNT(*) FROM comment WHERE review = 0")
	var cnt int
	if err = row.Scan(&cnt); nil != err {
		return 0, err
	}
	return cnt, nil
}

func modelCommentGetAllUnreviewed(page int, limit int) ([]*CommentModel, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}

	args := make([]interface{}, 0, 3)
	args = append(args)
	sqlExpr := `SELECT 
				id, 
				uid,
				reply_user, 
				is_sub, 
				sub_ref_id, 
				sub_to_uid,
				sub_to_user, 
				comment, 
				create_time, 
				update_time,
				agree,
				oppose,
				review FROM comment WHERE review = 0 ORDER BY create_time `
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

	replys := make([]*CommentModel, 0, 32)
	for rows.Next() {
		var reply CommentModel

		if err = rows.Scan(&reply.Id,
			&reply.Uid,
			&reply.ReplyUser,
			&reply.IsSub,
			&reply.SubRefId,
			&reply.SubToUid,
			&reply.SubToUser,
			&reply.Comment,
			&reply.CreateTime,
			&reply.UpdateTime,
			&reply.Agree,
			&reply.Oppose,
			&reply.Review); nil != err {
			return nil, err
		}
		replys = append(replys, &reply)
	}

	return replys, nil
}

func modelCommentGetArticleReply(uri string, page int, limit int, all bool) ([]*CommentModel, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}

	args := make([]interface{}, 0, 3)
	args = append(args, uri)
	sqlExpr := `SELECT 
				id, 
				uid,
				reply_user, 
				is_sub, 
				sub_ref_id, 
				sub_to_uid,
				sub_to_user, 
				comment, 
				create_time, 
				update_time,
				agree,
				oppose,
				review FROM comment WHERE uri = ?`
	if !all {
		sqlExpr += " AND review = 1 "
	}
	sqlExpr += " ORDER BY create_time "
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

	replys := make([]*CommentModel, 0, 32)
	for rows.Next() {
		var reply CommentModel

		if err = rows.Scan(&reply.Id,
			&reply.Uid,
			&reply.ReplyUser,
			&reply.IsSub,
			&reply.SubRefId,
			&reply.SubToUid,
			&reply.SubToUser,
			&reply.Comment,
			&reply.CreateTime,
			&reply.UpdateTime,
			&reply.Agree,
			&reply.Oppose,
			&reply.Review); nil != err {
			return nil, err
		}
		replys = append(replys, &reply)
	}

	return replys, nil
}

func modelNewComment(uri string, user *WebUser, comment string, parentId int, parentSubUser uint32) (int64, error) {
	db, err := getRawDB()
	if nil != err {
		return 0, err
	}

	if len(comment) > maxCommentLength {
		return 0, errors.New("Comment max length is 512 characters")
	}
	if len(uri) == 0 {
		return 0, errors.New("Invalid url for reply")
	}

	var reply CommentModel
	reply.Comment = comment
	reply.CreateTime = time.Now().Unix()
	reply.IsSub = false
	if 0 != parentId {
		reply.IsSub = true
		// Need get the parent comment
		pcomment, err := modelCommentGet(parentId)
		if nil != err {
			return 0, errors.New("Can't find parent comment")
		}
		if 0 != parentSubUser {
			subUser := modelWebUserGetUserByUid(parentSubUser)
			if nil == subUser {
				return 0, errors.New("Can't find parent sub user")
			}
			reply.SubToUid = subUser.Uid
			reply.SubToUser = subUser.UserName
		}
		reply.SubRefId = pcomment.Id
	}
	reply.Uid = user.Uid
	reply.ReplyUser = user.UserName
	reply.Uri = uri

	review := 0
	if user.Permission == kPermission_SuperAdmin {
		review = 1
	}

	if ret, err := db.Exec(`INSERT INTO comment (
										uid, 
										reply_user, 
										is_sub,
										sub_ref_id,
										sub_to_uid,
										sub_to_user, 
										uri, 
										comment, 
										create_time,
										agree,
										oppose,
										review) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		reply.Uid,
		reply.ReplyUser,
		reply.IsSub,
		reply.SubRefId,
		reply.SubToUid,
		reply.SubToUser,
		reply.Uri,
		reply.Comment,
		reply.CreateTime,
		reply.Agree,
		reply.Oppose,
		review); nil != err {
		return 0, err
	} else {
		return ret.LastInsertId()
	}
}

func modelCommentMarkReview(rid int) error {
	db, err := getRawDB()
	if nil != err {
		return err
	}

	_, err = db.Exec("UPDATE comment SET review = 1 WHERE id = ?", rid)
	return err
}

func modelCommentDelete(rid int) error {
	db, err := getRawDB()
	if nil != err {
		return err
	}

	_, err = db.Exec("DELETE FROM comment WHERE id = ?", rid)
	return err
}

func modelCommentDeleteByURI(uri string) error {
	db, err := getRawDB()
	if nil != err {
		return err
	}

	_, err = db.Exec("DELETE FROM comment WHERE uri = ?", uri)
	return err
}

func modelCommentDeleteSubRefID(srid int) error {
	db, err := getRawDB()
	if nil != err {
		return err
	}

	_, err = db.Exec("DELETE FROM comment WHERE sub_ref_id = ?", srid)
	return err
}
