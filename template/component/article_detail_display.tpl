{{define "article_detail_display"}}
<div id="article-container" class="media">
	<a class="pull-left" href="/project/{{.ProjectId}}/article/{{.Id}}" target="_blank">
		<img class="media-object" src="
		{{if eq .CoverImage ""}}
		{{getImagePath "/article_cover.png"}}
		{{else}}
		{{getImagePath .CoverImage}}
		{{end}}"
		width="219" height="148" style="border:1px solid #EEEEEE;padding:2px 2px 2px 2px;" />
	</a>
	<div class="media-body category-body">
		<h4>
			<a href="/project/{{.ProjectId}}/article/{{.Id}}" target="_blank">{{.ArticleTitle}}</a>
			{{if eq .Top 1}}
			<span style="color:red;margin-left:10px;">[置顶]</span>
			{{end}}
		</h4>
		<div class="category-base-info" style="padding-left:5px;border-left:2px solid #3472ef;">
			<i class="fa fa-smile-o"></i> 作者：<a href="/member/{{.ArticleAuthor}}" target="_blank">{{.ArticleAuthor}}</a>
			<span style="margin-left:8px;"><i class="fa fa-clock-o"></i> 更新时间：{{getTimeGapString .ActiveTime}}<span>
			<span style="margin-left:8px;"><i class="fa fa-server"></i> 分类：<a href="/project/{{.ProjectId}}/page/1" target="_blank">{{.ProjectName}}</a><span>
		</div>
		<div class="category-short-content">
			{{getThumb .ArticleContentHtml 60}}
		</div>
		<div style="border-bottom:1px solid #EEEEEE;margin-top:8px;margin-bottom:5px;"></div>
		<div class="category-ex-info">
			<i class="fa fa-hand-o-up"></i> 点击数：{{.Click}}
			<span style="margin-left:8px;"><i class="fa fa-reply"></i> 回复数：<span id="id-article-last-reply-{{.Id}}" class="article-last-reply" articleId="{{.Id}}">0</span><span>
		</div>
	</div>
</div>
{{end}}