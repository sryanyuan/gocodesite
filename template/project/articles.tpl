{{define "Title"}}文章{{end}}
{{define "importcss"}}
<link href="/static/css/articles.css" rel="stylesheet" />
{{end}}
{{define "importjs"}}
<!--script src="/static/js/project_category.js"></script-->
{{end}}
{{define "content"}}
<div id="id-content" class="container">
	<div class="row">
		<div class="col-md-6 col-md-offset-3">
			{{if gt .user.Permission 3}}
			<!--Administrator panel-->
			<p>
				<a href="/project/{{.project}}/new_article"><button type="button" class="btn btn-sm btn-primary">添加文章</button></a>
			</p>
			{{end}}
			<div class="articles-container">
				<dl>
				{{range .articles}}
					<dd>
						<a href="/member/{{.ArticleAuthor}}" class="pull-left" style="margin-right:10px;">
							<img class="img-rounded" src="{{$.imgPrefix}}/{{getMemberAvatar .ArticleAuthor}}" width="45" height="45" alt="{{.ArticleAuthor}}" >
						</a>
						<a href="/project/{{.ProjectName}}/article/{{.Id}}" class="title">
							{{.ArticleTitle}}
						</a>
						<div class="space"></div>
						<div class="info" style="margin-left:55px">
							<!--name-->
							<a href="/member/{{.ArticleAuthor}}">
								<strong>{{.ArticleAuthor}}</strong>
							</a>
							• {{getTimeGapString .ActiveTime}}
							{{if ne .ReplyAuthor ""}}
								• 最后回复来自 <a href="/member/{{.ReplyAuthor}}">{{.ReplyAuthor}}</a>
							{{end}}
						</div>
					</dd>
				{{end}}
				</dl>
			</div>
		</div>
	</div>
</div>
{{end}}