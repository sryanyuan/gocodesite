package gocodecc

import (
	"fmt"
	"net/http"
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
	registerRouter("/api/article", kPermission_Guest, apiArticlesGet, []string{http.MethodGet})
	registerRouter("/api/article/{articleId}", kPermission_Guest, apiArticleGet, []string{http.MethodGet})
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
	cnt, err := modelReplyGetCountByURI(fmt.Sprintf("/project/%d/article/%d", article.CategoryID, article.ArticleID))
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
	rsp.Content, err = convertMarkdown2HTML(article.ArticleContentMarkdown, summary)
	if nil != err {
		ctx.WriteAPIRspBadInternalError(err.Error())
		return
	}
	if ctx.config.CommentProvider == "native" {
		if err = fillArticleReplyCount(&rsp); nil != err {
			ctx.WriteAPIRspBadInternalError(err.Error())
			return
		}
	}

	ctx.WriteAPIRspOKWithMessage(&rsp)
}
