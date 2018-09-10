package gocodecc

import (
	"time"

	"github.com/cihub/seelog"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

type apiArticleCommentRsp struct {
	Id      int                     `json:"id"`
	Uid     int                     `json:"uid"`
	Name    string                  `json:"name"`
	Content string                  `json:"content"`
	Time    string                  `json:"time"`
	Agree   int                     `json:"agree"`
	Oppose  int                     `json:"oppose"`
	ToUid   int                     `json:"toUid"`
	ToUser  string                  `json:"toUser"`
	Review  int                     `json:"review"`
	Subs    []*apiArticleCommentRsp `json:"subs"`
}

type apiArticleCommentsRsp struct {
	Replys []*apiArticleCommentRsp `json:"replys"`
}

func getCommentsMergedByURI(uri string, user *WebUser) (*apiArticleCommentsRsp, error) {
	comments, err := modelCommentGetArticleReply(uri, 0, 0, user.Permission == kPermission_SuperAdmin)
	if nil != err {
		return nil, err
	}
	var rsp apiArticleCommentsRsp
	rsp.Replys = make([]*apiArticleCommentRsp, 0, len(comments))
	commentMap := make(map[int]*apiArticleCommentRsp)
	// Merge comments
	for _, comment := range comments {
		if comment.SubRefId == 0 {
			// Top comment
			var topComment apiArticleCommentRsp
			commentMap[comment.Id] = &topComment
			topComment.Id = comment.Id
			topComment.Uid = int(comment.Uid)
			tm := time.Unix(comment.CreateTime, 0)
			topComment.Time = tm.Format("2006-01-02 15:04:05")
			topComment.Content = comment.Comment
			topComment.Agree = comment.Agree
			topComment.Oppose = comment.Oppose
			topComment.Subs = make([]*apiArticleCommentRsp, 0, 32)
			topComment.Name = comment.ReplyUser
			topComment.Review = comment.Review
			rsp.Replys = append(rsp.Replys, &topComment)
		}
	}
	// Merge sub comments
	for _, comment := range comments {
		if comment.SubRefId == 0 {
			continue
		}
		topComment, ok := commentMap[comment.SubRefId]
		if !ok || nil == topComment {
			seelog.Errorf("Can't find parent comment while finding sub comment %d 's parent", comment.SubRefId)
			continue
		}
		var subComment apiArticleCommentRsp
		subComment.Id = comment.Id
		subComment.Uid = int(comment.Uid)
		subComment.Name = comment.ReplyUser
		tm := time.Unix(comment.CreateTime, 0)
		subComment.Time = tm.Format("2006-01-02 15:04:05")
		subComment.Content = comment.Comment
		subComment.Agree = comment.Agree
		subComment.Oppose = comment.Oppose
		subComment.ToUid = int(comment.SubToUid)
		subComment.ToUser = comment.SubToUser
		subComment.Review = comment.Review
		topComment.Subs = append(topComment.Subs, &subComment)
	}
	return &rsp, nil
}

func getCommentMergedByURI(uri string, commentId int, user *WebUser) (*apiArticleCommentRsp, error) {
	var rsp apiArticleCommentRsp
	// Find parent first
	topComment, err := modelCommentGet(commentId)
	if nil != err {
		return nil, err
	}
	if nil != topComment {
		rsp.Id = topComment.Id
		rsp.Uid = int(topComment.Uid)
		tm := time.Unix(topComment.CreateTime, 0)
		rsp.Time = tm.Format("2006-01-02 15:04:05")
		rsp.Content = topComment.Comment
		rsp.Agree = topComment.Agree
		rsp.Oppose = topComment.Oppose
		rsp.Name = topComment.ReplyUser
		rsp.Review = topComment.Review
		// Get subs
		subs, err := modelCommentGetSubs(uri, commentId, 0, 0, user.Permission == kPermission_SuperAdmin)
		if nil != err {
			return nil, err
		}
		if nil != subs {
			rsp.Subs = make([]*apiArticleCommentRsp, 0, len(subs))
			for _, comment := range subs {
				// Top comment
				var sub apiArticleCommentRsp
				sub.Id = comment.Id
				sub.Uid = int(comment.Uid)
				tm := time.Unix(comment.CreateTime, 0)
				sub.Time = tm.Format("2006-01-02 15:04:05")
				sub.Content = comment.Comment
				sub.Agree = comment.Agree
				sub.Oppose = comment.Oppose
				sub.Name = comment.ReplyUser
				sub.ToUser = comment.SubToUser
				sub.ToUid = int(comment.SubToUid)
				sub.Review = comment.Review
				rsp.Subs = append(rsp.Subs, &sub)
			}
		}
	}
	return &rsp, nil
}
