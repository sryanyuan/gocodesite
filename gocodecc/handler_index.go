package gocodecc

import "fmt"

var homeRenderTpls []string = []string{
	"template/component/article_detail_display.html",
	"template/home.html",
}

func indexHandler(ctx *RequestContext) {
	// Get top articles
	topArticles, err := modelProjectArticleGetAllTopArticles(0, 0)
	if nil != err {
		ctx.RenderMessagePage("错误", err.Error(), false)
		return
	}
	//	get recent articles
	recentArticles, err := modelProjectArticleGetRecentNotTopArticles(0, 5)
	if nil != err {
		ctx.RenderMessagePage("错误", err.Error(), false)
		return
	}

	//	get article count
	articleCount, err := modelProjectArticleGetArticleCountAll()
	if nil != err {
		ctx.RenderMessagePage("错误", err.Error(), false)
		return
	}

	//	get category
	category, err := modelProjectCategoryGetAllSimple()
	if nil != err {
		ctx.RenderMessagePage("错误", err.Error(), false)
		return
	}

	//	get member count
	memberCount, err := modelWebUserGetCount()
	if nil != err {
		ctx.RenderMessagePage("错误", err.Error(), false)
		return
	}

	// Get reply count
	replyCount, err := modelReplyGetCount()
	if nil != err {
		ctx.RenderMessagePage("错误", err.Error(), false)
		return
	}

	// Handle private article
	for _, article := range recentArticles {
		articleApplyPrivate(ctx.user, article)
	}
	for _, article := range topArticles {
		articleApplyPrivate(ctx.user, article)
	}

	dataCtx := map[string]interface{}{
		"active":         "home",
		"topArticles":    topArticles,
		"recentArticles": recentArticles,
		"articleCount":   articleCount,
		"category":       category,
		"memberCount":    memberCount,
		"createSiteTime": metaInfoCreateSiteTime,
		"replyCount":     replyCount,
	}
	if ctx.config.CommentProvider == "native" {
		// Get all comment count
		for _, v := range topArticles {
			cnt, err := modelReplyGetCountByURI(fmt.Sprintf("/project/%d/article/%d", v.ProjectId, v.Id))
			if nil == err {
				v.ReplyCount = cnt
			}
		}
		for _, v := range recentArticles {
			cnt, err := modelReplyGetCountByURI(fmt.Sprintf("/project/%d/article/%d", v.ProjectId, v.Id))
			if nil == err {
				v.ReplyCount = cnt
			}
		}
	}
	dataHtml := renderTemplate(ctx, homeRenderTpls, dataCtx)
	ctx.w.Write(dataHtml)
}
