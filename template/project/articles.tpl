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
		<div class="col-md-8 col-md-offset-2">
			{{if gt .user.Permission 3}}
			<!--Administrator panel-->
			<p>
				<a href="/project/{{.project}}/cmd/new_article"><button type="button" class="btn btn-sm btn-success">添加文章</button></a>
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
						项目
					</a>
				</li>
			</div>
			<div class="articles-container">
				<dl>
				{{range .articles}}
					<dd>
						<a href="/member/{{.ArticleAuthor}}" class="pull-left" style="margin-right:10px;">
							<img class="img-rounded" src="{{$.imgPrefix}}/{{getMemberAvatar .ArticleAuthor}}" width="45" height="45" alt="{{.ArticleAuthor}}" >
						</a>
						<a href="/project/{{.ProjectName}}/article/{{.Id}}" class="title">
							{{.ArticleTitle}}
							{{if eq .Top 1}}
							<i class="fa fa-angle-up" style="margin-left:25px;"></i>
							{{end}}
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
			<div style="text-align:center;">
				<nav>
					<ul class="pagination">
						<!--have more than 1 page-->
						{{if gt .pages 1}}
						<!--previous page-->
						{{if gt .page 1}}
						<li><a href="/project/{{.project}}/page/{{minusInt .page}}" aria-lable="Previous"><span aria-hidden="true">&laquo;</span></a></li>
						{{else}}
						<li class="disabled"><a href="javascript:void(0);" aria-lable="Previous"><span aria-hidden="true">&laquo;</span></a></li>
						{{end}}
						<!--fill pages-->
						{{$pageRange := getPageRange .page .showPages}}
						{{range $i, $v := $pageRange}}
							{{if eq $v $.page}}
							<li class="active"><a href="javascript:void(0);">{{$v}}</a></li>
							{{else}}
								{{if gt $v $.pages}}
									<li class="disabled"><a href="javascript:void(0);">{{$v}}</a></li>
								{{else}}
									<li><a href="/project/{{$.project}}/page/{{$v}}">{{$v}}</a></li>
								{{end}}
							{{end}}
						{{end}}
						<!--next page-->
						{{if lt .page .pages}}
						<li><a href="/project/{{.project}}/page/{{addInt .page}}" aria-lable="Next"><span aria-hidden="true">&raquo;</span></a></li>
						{{else}}
						<li class="disabled"><a href="javascript:void(0);" aria-lable="Next"><span aria-hidden="true">&raquo;</span></a></li>
						{{end}}
						
						{{end}}
					</ul>
				</nav>
			</div>
		</div>
	</div>
</div>
{{end}}