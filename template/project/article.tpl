{{define "Title"}}文章{{end}}
{{define "importcss"}}
<link href="/static/css/editormd.min.css" rel="stylesheet" />
<link href="/static/css/article.css" rel="stylesheet" />
{{end}}
{{define "importjs"}}
<script src="/static/js/editormd.min.js"></script>
<script src="/static/js/article.js"></script>
{{end}}
{{define "content"}}
<div id="id-content" class="container">
	<div class="row">
		<div class="col-md-8">
			<!--Administrator panel-->
			{{if gt .user.Permission 3}}
			<p>
				<a href="/project/{{.article.ProjectName}}/cmd/edit_article?articleId={{.article.Id}}"><button type="button" class="btn btn-sm btn-success">编辑文章</button></a>
				<button id="id-article-del" type="button" onclick="deleteArticle(this, {{.article.Id}})" class="btn btn-sm btn-danger">删除文章</button>
			</p>
			<hr/>
			{{else if eq .user.NickName .article.ArticleAuthor}}
			<p>
				<button id="id-article-edit" type="button" class="btn btn-sm btn-success">编辑文章</button>
			</p>
			<hr/>
			{{end}}
			<div class="breadcrumb">
				<li>
					<a href="/">
						<i class="fa fa-home"></i>首页
					</a>
				</li>
				<li>
					<a href="/project/{{.article.ProjectName}}/page/1">
						{{.article.ProjectName}}
					</a>
				</li>
			</div>
			<div class="content">
				<div class="page-header">
					<div style="text-align:center;"><h1>{{.article.ArticleTitle}}</h1></div>
					<div class="btn-group btn-group-sm pull-right">
						<a href="/ajax/article_like/{{.article.Id}}" class="btn btn-default" title="{{.article.Like}} 赞">
							<i class="fa fa-heart"></i>
						</a>
						<a href="/ajax/article_fav/{{.article.Id}}" class="btn btn-default" title="收藏">
							<i class="fa fa-star"></i>
						</a>
					</div>
					<div style="padding-bottom:5px;">
						<span style="margin-bottom:5px;">由 <a href="/member/{{.article.ArticleAuthor}}">{{.article.ArticleAuthor}}</a> 在 {{getTimeGapString .article.PostTime}} 发布 {{.article.Click}} 次点击</span>
					</div>
				</div>
				<div style="padding-top:10px;" class="body editormd-preview-container">
					{{.article.ArticleContentMarkdown}}
				</div>
			</div>
		</div>
		<div class="col-md-4">
			<div class="panel panel-default">
				<div class="panel-heading">
					<h3 class="panel-title">作者</h3>
				</div>
				<div class="panel-body">
					<div>
						<a href="/member/{{.article.ArticleAuthor}}">
							<img class="gravatar img-rounded" style="float:left;margin-right:10px;" src="{{.imgPrefix}}/{{getMemberAvatar .article.ArticleAuthor}}" width="42" height="42"></img>
						</a>
						<h4>
							<a href="/member/{{.article.ArticleAuthor}}">{{.article.ArticleAuthor}}</a>
							<br/>
							<small>{{.author.Mood}}</small>
						</h4>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
{{end}}