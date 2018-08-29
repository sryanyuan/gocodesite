package gocodecc

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cihub/seelog"
	"github.com/dchest/captcha"
)

const (
	// Post time with limit articles show in index page
	articlesGetModePostTime = iota
	// Get articles in archive time mode, support limit and page
	articlesGetModeArchive
	// Get articles with category, support limit and page
	articlesGetModeCategory
	articlesGetModeTotal
)

func init() {
	registerApi("/api/article", kPermission_Guest, apiArticlesGet, []string{http.MethodGet})
	registerApi("/api/article", kPermission_Guest, apiArticlePost, []string{http.MethodPost})
	registerApi("/api/article/{articleId}", kPermission_Guest, apiArticleGet, []string{http.MethodGet})
	registerApi("/api/article/{articleId}", kPermission_SuperAdmin, apiArticlePut, []string{http.MethodPut})
	registerApi("/api/article/{articleId}", kPermission_SuperAdmin, apiArticleDelete, []string{http.MethodDelete})
	registerApi("/api/article/{articleId}/comment", kPermission_Guest, apiArticleCommentsGet, []string{http.MethodGet})
	registerApi("/api/article/{articleId}/comment/{commentId}", kPermission_Guest, apiArticleCommentGet, []string{http.MethodGet})
	registerApi("/api/article/{articleId}/comment", kPermission_User, apiArticleCommentPost, []string{http.MethodPost})
	registerApi("/api/comment/{commentId}", kPermission_SuperAdmin, apiArticleCommentDelete, []string{http.MethodDelete})
}

type apiArticleRsp struct {
	AuthorID     int    `json:"authorId"`
	AuthorName   string `json:"authorName"`
	Top          bool   `json:"top"`
	Category     string `json:"category"`
	CategoryID   int    `json:"categoryId"`
	ArticleID    int    `json:"articleId"`
	Content      string `json:"content"`
	Title        string `json:"title"`
	PostDatetime string `json:"postDatetime"`
	ReplyCount   int    `json:"replyCount"`
}

type apiArticlesRsp struct {
	Articles []*apiArticleRsp `json:"articles"`
	Total    int              `json:"total"`
	Pages    int              `json:"pages"`
}

func fillArticleReplyCount(article *apiArticleRsp) error {
	// Get all comment count
	cnt, err := modelCommentGetTopCount(fmt.Sprintf("article:%d", article.ArticleID))
	if nil == err {
		article.ReplyCount = cnt
	}
	return err
}

func fillArticlesReplyCount(articles *apiArticlesRsp) error {
	for _, v := range articles.Articles {
		if err := fillArticleReplyCount(v); nil != err {
			return err
		}
	}
	return nil
}

func apiArticlesGet(ctx *RequestContext) {
	mode := ctx.GetFormValueInt("mode", 0)
	switch mode {
	case articlesGetModePostTime:
		{
			limit := ctx.GetFormValueInt("limit", 0)
			if limit <= 0 {
				limit = 10
			}
			topArticles, err := modelProjectArticleGetAllTopArticles(0, 0)
			if nil != err {
				ctx.WriteAPIRspBadInternalError(err.Error())
				return
			}
			recentArticles, err := modelProjectArticleGetRecentNotTopArticles(0, limit)
			if nil != err {
				ctx.WriteAPIRspBadInternalError(err.Error())
				return
			}
			var rsp apiArticlesRsp
			rsp.Articles = make([]*apiArticleRsp, 0, len(topArticles)+len(recentArticles))
			for _, v := range topArticles {
				var item apiArticleRsp
				item.ArticleID = v.Id
				item.Category = v.ProjectName
				item.CategoryID = v.ProjectId
				item.PostDatetime = tplfn_getTimeGapString(v.PostTime)
				item.AuthorName = v.ArticleAuthor
				if author := modelWebUserGetUserByUserName(v.ArticleAuthor); nil != author {
					item.AuthorID = int(author.Uid)
				}
				item.Title = v.ArticleTitle
				if v.Top != 0 {
					item.Top = true
				}
				rsp.Articles = append(rsp.Articles, &item)
			}
			for _, v := range recentArticles {
				var item apiArticleRsp
				item.ArticleID = v.Id
				item.Category = v.ProjectName
				item.CategoryID = v.ProjectId
				item.PostDatetime = tplfn_getTimeGapString(v.PostTime)
				item.AuthorName = v.ArticleAuthor
				if author := modelWebUserGetUserByUserName(v.ArticleAuthor); nil != author {
					item.AuthorID = int(author.Uid)
				}
				item.Title = v.ArticleTitle
				if v.Top != 0 {
					item.Top = true
				}
				rsp.Articles = append(rsp.Articles, &item)
			}
			if ctx.config.CommentProvider == "native" {
				if err = fillArticlesReplyCount(&rsp); nil != err {
					ctx.WriteAPIRspBadInternalError(err.Error())
					return
				}
			}
			ctx.WriteAPIRspOKWithMessage(&rsp)
		}
	case articlesGetModeCategory:
		{
			page := ctx.GetFormValueInt("page", 0)
			limit := ctx.GetFormValueInt("limit", 10)
			category := ctx.GetFormValueInt("category", 0)
			articles, pages, err := modelProjectArticleGetArticles(category, page, limit)
			if nil != err {
				ctx.WriteAPIRspBadInternalError(err.Error())
				return
			}
			var rsp apiArticlesRsp
			rsp.Pages = pages
			rsp.Articles = make([]*apiArticleRsp, 0, len(articles))
			for _, v := range articles {
				var item apiArticleRsp
				item.ArticleID = v.Id
				item.Category = v.ProjectName
				item.CategoryID = v.ProjectId
				item.PostDatetime = tplfn_getTimeGapString(v.PostTime)
				item.AuthorName = v.ArticleAuthor
				if author := modelWebUserGetUserByUserName(v.ArticleAuthor); nil != author {
					item.AuthorID = int(author.Uid)
				}
				item.Title = v.ArticleTitle
				if v.Top != 0 {
					item.Top = true
				}
				rsp.Articles = append(rsp.Articles, &item)
			}
			if ctx.config.CommentProvider == "native" {
				if err = fillArticlesReplyCount(&rsp); nil != err {
					ctx.WriteAPIRspBadInternalError(err.Error())
					return
				}
			}
			ctx.WriteAPIRspOKWithMessage(&rsp)
		}
	default:
		{
			ctx.WriteAPIRspBadInternalError("invalid mode")
		}
	}
}

type apiArticlePostArg struct {
	CategoryId int    `json:"category"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CoverImage string `json:"coverImage"`
}

type apiArticlePostRsp struct {
	ArticleId int `json:"articleId"`
}

func apiArticlePost(ctx *RequestContext) {
	var arg apiArticlePostArg
	if err := ctx.readFromBody(&arg); nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}

	var prj ProjectCategoryItem
	if err := modelProjectCategoryGetByProjectId(arg.CategoryId, &prj); nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}

	if len(arg.Title) == 0 || len(arg.Title) > kArticleTitleLimit {
		ctx.WriteAPIRspBadInternalError("title length out of range")
		return
	}
	if len(arg.Content) == 0 || len(arg.Content) > kArticleContentLimit {
		ctx.WriteAPIRspBadInternalError("title content out of range")
		return
	}

	var postArticle ProjectArticleItem
	postArticle.ActiveTime = time.Now().Unix()
	postArticle.PostTime = time.Now().Unix()
	postArticle.ArticleTitle = arg.Title
	postArticle.ArticleAuthor = ctx.user.NickName
	postArticle.ArticleContentMarkdown = arg.Content
	postArticle.ArticleContentHtml, _ = convertMarkdown2HTML(arg.Content, 0)
	postArticle.ProjectName = prj.ProjectName
	postArticle.ProjectId = prj.Id
	postArticle.CoverImage = arg.CoverImage
	articleId, err := modelProjectArticleNewArticle(&postArticle)
	if nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	var rsp apiArticlePostRsp
	rsp.ArticleId = int(articleId)
	ctx.WriteAPIRspOKWithMessage(&rsp)
}

func apiArticleDelete(ctx *RequestContext) {
	articleId := ctx.GetURLVarInt64("articleId", 0)
	if 0 == articleId {
		ctx.WriteAPIRspBadInternalError("invalid article id")
		return
	}
	// Find article
	article, err := modelProjectArticleGet(int(articleId))
	if nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	if err = modelProjectArticleDelete(int(articleId), article.ProjectId); nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	ctx.WriteAPIRspOK(nil)
}

type apiArticlePutArg struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func apiArticlePut(ctx *RequestContext) {
	articleId := ctx.GetURLVarInt64("articleId", 0)
	if 0 == articleId {
		ctx.WriteAPIRspBadInternalError("invalid article id")
		return
	}
	var arg apiArticlePutArg
	if err := ctx.readFromBody(&arg); nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	if len(arg.Content) == 0 || len(arg.Content) > kArticleContentLimit {
		ctx.WriteAPIRspBadInternalError("Content length out of range")
		return
	}
	if len(arg.Title) == 0 || len(arg.Title) > kArticleTitleLimit {
		ctx.WriteAPIRspBadInternalError("Title length out of range")
		return
	}
	article, err := modelProjectArticleGet(int(articleId))
	if nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	colsEdit := []string{"active_time", "edit_time"}
	article.ActiveTime = time.Now().Unix()
	article.EditTime = time.Now().Unix()
	if article.ArticleTitle != arg.Title {
		article.ArticleTitle = arg.Title
		colsEdit = append(colsEdit, "article_title")
	}
	if article.ArticleContentMarkdown != arg.Content {
		article.ArticleContentHtml, err = convertMarkdown2HTML(arg.Content, 0)
		article.ArticleContentMarkdown = arg.Content
		colsEdit = append(colsEdit, "article_content_html")
		colsEdit = append(colsEdit, "article_content_markdown")
	}
	_, err = modelProjectArticleEditArticle(article, colsEdit)
	if nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	ctx.WriteAPIRspOK(nil)
}

func apiArticleGet(ctx *RequestContext) {
	articleId := ctx.GetURLVarInt64("articleId", 0)
	if 0 == articleId {
		ctx.WriteAPIRspBadInternalError("invalid article id")
		return
	}
	summary := ctx.GetFormValueInt("summary", 0)
	article, err := modelProjectArticleGet(int(articleId))
	if nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	mk := ctx.GetFormValueInt("mk", 0)

	var rsp apiArticleRsp
	rsp.ArticleID = article.Id
	rsp.CategoryID = article.ProjectId
	rsp.Category = article.ProjectName
	rsp.Title = article.ArticleTitle
	rsp.PostDatetime = tplfn_getTimeGapString(article.PostTime)
	if article.Top != 0 {
		rsp.Top = true
	}
	rsp.AuthorName = article.ArticleAuthor
	if author := modelWebUserGetUserByUserName(article.ArticleAuthor); nil != author {
		rsp.AuthorID = int(author.Uid)
	}
	// Convert markdown to html
	if mk == 0 {
		rsp.Content, err = convertMarkdown2HTML(article.ArticleContentMarkdown, summary)
		if nil != err {
			ctx.WriteAPIRspBadInternalError(err.Error())
			return
		}
	} else {
		rsp.Content = article.ArticleContentMarkdown
	}

	if ctx.config.CommentProvider == "native" {
		if err = fillArticleReplyCount(&rsp); nil != err {
			ctx.WriteAPIRspBadInternalError(err.Error())
			return
		}
	}

	ctx.WriteAPIRspOKWithMessage(&rsp)
}

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
	Subs    []*apiArticleCommentRsp `json:"subs"`
}

type apiArticleCommentsRsp struct {
	Replys []*apiArticleCommentRsp `json:"replys"`
}

func apiArticleCommentsGet(ctx *RequestContext) {
	uri := fmt.Sprintf("article:%s", ctx.GetURLVarString("articleId"))
	comments, err := modelCommentGetArticleReply(uri, 0, 0)
	if nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
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
		topComment.Subs = append(topComment.Subs, &subComment)
	}
	ctx.WriteAPIRspOKWithMessage(&rsp)
}

func apiArticleCommentGet(ctx *RequestContext) {
	uri := fmt.Sprintf("article:%s", ctx.GetURLVarString("articleId"))
	commentId := int(ctx.GetURLVarInt64("commentId", 0))
	var rsp apiArticleCommentRsp
	// Find parent first
	topComment, err := modelCommentGet(commentId)
	if nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
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
		// Get subs
		subs, err := modelCommentGetSubs(uri, commentId, 0, 0)
		if nil != err {
			ctx.WriteAPIRspBadInternalError(err.Error())
			return
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
				rsp.Subs = append(rsp.Subs, &sub)
			}
		}
	}
	ctx.WriteAPIRspOKWithMessage(&rsp)
}

type apiArticleCommentPostArg struct {
	Content   string `json:"content"`
	URI       string `json:"uri"`
	SubRefID  int    `json:"subRefID"`
	ToUser    uint32 `json:"toUser"`
	CaptchaId string `json:"captchaId"`
	Solution  string `json:"solution"`
}

func apiArticleCommentPost(ctx *RequestContext) {
	var arg apiArticleCommentPostArg
	if err := ctx.readFromBody(&arg); nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	if arg.Solution == "" || arg.CaptchaId == "" {
		ctx.WriteAPIRspBadInternalError("invalid captcha input")
		return
	}
	if !captcha.VerifyString(arg.CaptchaId, arg.Solution) {
		ctx.WriteAPIRspBadInternalError("invalid catpcha")
		return
	}
	if len(arg.Content) < 5 || len(arg.Content) > 128 {
		ctx.WriteAPIRspBadInternalError("content is too long")
		return
	}
	// Check parent comment has same uri
	if 0 != arg.SubRefID {
		parentComment, err := modelCommentGet(arg.SubRefID)
		if nil != err {
			ctx.WriteAPIRspBadInternalError(err.Error())
			return
		}
		if parentComment.IsSub {
			ctx.WriteAPIRspBadInternalError("Can't reply to sub reply")
			return
		}
		if parentComment.Uri != arg.URI {
			ctx.WriteAPIRspBadInternalError("parent uri not equal")
			return
		}
	}

	if _, err := modelNewComment(arg.URI, ctx.user, arg.Content, arg.SubRefID, arg.ToUser); nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	ctx.WriteAPIRspOK(nil)
}

func apiArticleCommentDelete(ctx *RequestContext) {
	commentId := ctx.GetURLVarInt64("commentId", 0)
	comment, err := modelCommentGet(int(commentId))
	if nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	if nil == comment {
		ctx.WriteAPIRspOK(nil)
		return
	}
	if err := modelCommentDelete(int(commentId)); nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	if err = modelCommentDeleteSubRefID(int(commentId)); nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}

	ctx.WriteAPIRspOK(nil)
}
