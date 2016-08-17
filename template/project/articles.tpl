{{define "Title"}}{{.category.ProjectName}}{{end}}
{{define "importcss"}}
<link href="/static/css/articles.css" rel="stylesheet" />
{{end}}
{{define "importjs"}}
<script src="/static/js/articles.js"></script>
{{end}}
{{define "content"}}
<div id="id-content" class="container">
	<div class="row">
		<div class="col-md-8 col-md-offset-2">
			{{if canPost .category .user}}
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
						分类
					</a>
				</li>
				<li>{{.category.ProjectName}}</li>
			</div>
			<div id="articles" articleCount="{{len .articles}}" class="articles-container">
				{{$articleCount := len .articles}}
				{{if eq $articleCount 0}}
				<h3 class="section-title-s1" style="max-width:none;">当前还有没有创建任何主题噢！</h3>
				{{else}}
				{{range .articles}}
				{{template "article_detail_display" .}}
				{{end}}
				{{end}}
			</div>
			
			<div style="text-align:center;">
				<nav>
					<ul class="pagination">
						<!--have more than 1 page-->
						{{if gt .pages 0}}
						<!--previous page-->
						{{if gt .page 1}}
						<li><a href="/project/{{.project}}/page/{{minusInt .page 1}}" aria-lable="Previous"><span aria-hidden="true">&laquo;</span></a></li>
						{{else}}
						<li class="disabled"><a href="javascript:void(0);" aria-lable="Previous" style="background-color:#F2F2F2"><span aria-hidden="true">&laquo;</span></a></li>
						{{end}}
						<!--fill pages-->
						{{$pageRange := getPageRange .page .showPages .pages}}
						{{range $i, $v := $pageRange}}
							<!--first page-->
							{{if eq $i 0}}
							{{if gt $v 2}}
							<li><a href="/project/{{$.project}}/page/1">1</a></li>
							<li class="disabled"><a href="javascript:void(0)" style="background-color:#F2F2F2">...</a></li>
							{{end}}
							{{end}}
							
							{{if eq $v $.page}}
							<li class="active"><a href="javascript:void(0);">{{$v}}</a></li>
							{{else}}
								{{if gt $v $.pages}}
									<li class="disabled"><a href="javascript:void(0);">{{$v}}</a></li>
								{{else}}
									<li><a href="/project/{{$.project}}/page/{{$v}}">{{$v}}</a></li>
								{{end}}
							{{end}}
							
							<!--last page-->
							{{$lastPageIndex := len $pageRange}}
							{{$lastPageIndex := minusInt $lastPageIndex 1}}
							{{$lastPage := minusInt $.pages 1}}
							{{if eq $i $lastPageIndex}}
							{{if lt $v $lastPage}}
							<li class="disabled"><a href="javascript:void(0)"  style="background-color:#F2F2F2">...</a></li>
							<li><a href="/project/{{$.project}}/page/{{$.pages}}">{{$.pages}}</a></li>
							{{end}}
							{{end}}
						{{end}}
						<!--next page-->
						{{if lt .page .pages}}
						<li><a href="/project/{{.project}}/page/{{addInt .page 1}}" aria-lable="Next"><span aria-hidden="true">&raquo;</span></a></li>
						{{else}}
						<li class="disabled"><a href="javascript:void(0);" aria-lable="Next" style="background-color:#F2F2F2"><span aria-hidden="true">&raquo;</span></a></li>
						{{end}}
						
						{{end}}
					</ul>
				</nav>
			</div>
		</div>
	</div>
</div>
{{end}}