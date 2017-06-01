package gocodecc

import (
	"strconv"
)

var moodRenderTpls = []string{
	"template/mood.html",
}

func moodHandler(ctx *RequestContext) {
	dataCtx := map[string]interface{}{
		"active": "mood",
	}

	ctx.r.ParseForm()
	pageStr := ctx.r.Form.Get("p")
	var page int
	var err error

	if len(pageStr) == 0 {
		page = 1
	} else {
		page, err = strconv.Atoi(ctx.r.Form.Get("p"))
		if nil != err {
			page = 1
		}
		if 0 >= page {
			page = 1
		}
	}

	//	get contents
	itemsPerPage := 25
	showPages := 5
	moods, err := modelMoodInfoGet(page-1, itemsPerPage)
	if nil != err {
		ctx.RenderMessagePage("错误", err.Error(), false)
		return
	}

	//	get pages
	pages, err := modelMoodGetCount()
	pages = (pages + itemsPerPage + 1) / itemsPerPage
	if nil != err {
		ctx.RenderMessagePage("错误", err.Error(), false)
		return
	}

	if page > pages {
		ctx.RenderMessagePage("错误", "Invalid page index", false)
		return
	}

	dataCtx["moods"] = moods
	dataCtx["pages"] = pages
	dataCtx["page"] = page
	dataCtx["showPages"] = showPages
	dataHtml := renderTemplate(ctx, moodRenderTpls, dataCtx)
	ctx.w.Write(dataHtml)
}
