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
		<div class="col-md-8 col-md-offset-0">
			<h2 class="section-title-s2">最近的文章</h2>
			<div id="recentArticles" class="articles-container">
				{{range .recentArticles}}
				{{template "article_detail_display" .}}
				{{end}}
			</div>
		</div>
		<div class="col-md-3 col-md-offset-0">
			<h2 class="section-title-s2"><a href="/project">主题目录</a></h2>
			<div class="section-category">
				<ul class="posts" style="list-style:none;">
					{{range .category}}
					<li class="post-item">
						<a href="/project/{{.Id}}/page/1">{{.ProjectName}}</a>
						<span style="float:right;">{{.ItemCount}}</span>
					</li>
					{{end}}
				</ul>
			</div>
			<div style="height:25px;"></div>
			<h2 class="section-title-s2">统计</h2>
			<div class="section-statistics">
				<p>主题数：{{.articleCount}}</p>
			</div>
		</div>
	</div>
</div>
{{end}}