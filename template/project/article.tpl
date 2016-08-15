{{define "Title"}}{{.article.ArticleTitle}}{{end}}
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
<!--Modal dialogs-->
	<div id="modalDeleteConfirm" class="modal fade in" role="dialog" aria-hidden="true" style="display: none;">
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<a class="close" data-dismiss="modal">×</a>
					<h4 id="id-modaldeleteconfirm-text"></h4>
				</div>
				<div class="modal-footer">
					<a href="#" onclick="submitDeleteArticle(this, '/ajax/article_delete')" class="btn btn-success">确定</a>
					<a href="#" class="btn" data-dismiss="modal">取消</a>
				</div>
			</div>
		</div>
	</div>
	<div id="modalAlert" class="modal fade in" role="dialog" aria-hidden="true" style="display: none;">
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<a class="close" data-dismiss="modal">×</a>
					<h3 style="color:#FE2E2E" id="id-modalalert-text">hint</h3>
				</div>
				<div class="modal-footer">
					<a href="#" class="btn" data-dismiss="modal">关闭</a>
				</div>
			</div>
		</div>
	</div>
	<div class="row">
		<div class="col-md-8">
			<!--Administrator panel-->
			{{if gt .user.Permission 3}}
			<p>
				<a href="/project/{{.article.ProjectId}}/cmd/edit_article?articleId={{.article.Id}}"><button type="button" class="btn btn-sm btn-success">编辑文章</button></a>
				{{if eq .article.Top 0}}
				<button id="id-article-top" type="button" articleId="{{.article.Id}}" articleTitle="{{.article.ArticleTitle}}" onclick="topArticle(this, true, {{.article.Id}})" class="btn btn-sm btn-success">置顶文章</button>
				{{else}}
				<button id="id-article-top" type="button" articleId="{{.article.Id}}" articleTitle="{{.article.ArticleTitle}}" onclick="topArticle(this, false, {{.article.Id}})" class="btn btn-sm btn-success">取消置顶</button>
				{{end}}
				<button id="id-article-del" type="button" articleId="{{.article.Id}}" articleTitle="{{.article.ArticleTitle}}" onclick="deleteArticle(this, {{.article.Id}})" class="btn btn-sm btn-danger">删除文章</button>
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
					<a href="/project">
						分类
					</a>
				</li>
				<li>
					<a href="/project/{{.article.ProjectId}}/page/1">
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
							{{.article.Like}}
						</a>
						<a href="/ajax/article_fav/{{.article.Id}}" class="btn btn-default" title="收藏">
							<i class="fa fa-star"></i>
						</a>
					</div>
					<div style="padding-bottom:5px;">
						<span style="margin-bottom:5px;">由 <a href="/member/{{.article.ArticleAuthor}}">{{.article.ArticleAuthor}}</a> 在 {{getTimeGapString .article.PostTime}} 发布 {{.article.Click}} 次点击</span>
					</div>
				</div>
				<div class="body editormd-preview-container">
					{{$content := .article.ArticleContentHtml}}
					{{convertToHtml $content}}
				</div>
			</div>
			{{template "comment_article_html" .}}
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