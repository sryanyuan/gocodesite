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

	dataCtx := map[string]interface{}{
		"active":         "home",
		"topArticles":    topArticles,
		"recentArticles": recentArticles,
		"articleCount":   articleCount,
		"category":       category,
		"memberCount":    memberCount,
		"createSiteTime": metaInfoCreateSiteTime,
	}
	if ctx.config.CommentProvider == "native" {
		// Get all comment count
		for _, v := range topArticles {
			cnt, err := modelReplyGetCount(fmt.Sprintf("/project/%d/article/%d", v.ProjectId, v.Id))
			if nil == err {
				v.ReplyCount = cnt
			}
		}
		for _, v := range recentArticles {
			cnt, err := modelReplyGetCount(fmt.Sprintf("/project/%d/article/%d", v.ProjectId, v.Id))
			if nil == err {
				v.ReplyCount = cnt
			}
		}
	}
	dataHtml := renderTemplate(ctx, homeRenderTpls, dataCtx)
	ctx.w.Write(dataHtml)
}