{{define "Title"}}sryan的个人小驿站 分享开发的过程与成果{{end}}
{{define "importcss"}}
<link href="/static/css/home.css" rel="stylesheet" />
<link href="/static/css/articles.css" rel="stylesheet" />
{{end}}
{{define "importjs"}}{{end}}
{{define "content"}}
<div id="id-content" class="container theme-showcase" role="main">
	<!--banner-->
	<div class="row">
		<div class="col-md-7 col-md-offset-1">
			<h2 class="section-title-s2">
				最近的文章
			</h2>
			<div id="recentArticles" class="articles-container">
				{{range .recentArticles}}
				{{template "article_detail_display" .}}
				{{end}}
			</div>
		</div>
	</div>
</div>
{{end}}