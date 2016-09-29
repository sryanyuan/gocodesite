{{define "Title"}}闲言碎语{{end}}
{{define "importcss"}}
<link href="/static/css/mood.css" rel="stylesheet" />
{{end}}
{{define "importjs"}}
<!--script src="/static/js/articles.js"></script-->
{{end}}
{{define "content"}}
<div id="id-content" class="container">
	<div class="row">
		<div class="col-md-10 col-md-offset-1">
			<div class="breadcrumb">
				<li>
					<a href="/">
						<i class="fa fa-home"></i>首页
					</a>
				</li>
				<li>闲言碎语</li>
			</div>
			<div id="moods" class="mood-container">
				{{$moodCount := len .moods}}
				{{if eq $moodCount 0}}
				<h3 class="section-title-s1" style="max-width:none;">当前还没有创建任何闲言碎语噢！</h3>
				{{else}}
				<div>
				{{range .moods}}
					<ul class="arrow-box">
						<div style="overflow:hidden;">
							<p>
								{{if ne .Image ""}}
								<img src="{{getMoodImagePath .Image}}"/>
								{{end}}
								{{.Mood}}
							</p>
						</div>
						<span class="dateview">{{formatDate .PostTime}}</span>
					</ul>
				{{end}}
				</div>
				{{end}}
			</div>
			
			<div style="text-align:center;">
				<nav>
					<ul class="pagination">
						<!--have more than 1 page-->
						{{if gt .pages 0}}
						<!--previous page-->
						{{if gt .page 1}}
						<li><a href="/mood?p={{minusInt .page 1}}" aria-lable="Previous"><span aria-hidden="true">&laquo;</span></a></li>
						{{else}}
						<li class="disabled"><a href="javascript:void(0);" aria-lable="Previous" style="background-color:#F2F2F2"><span aria-hidden="true">&laquo;</span></a></li>
						{{end}}
						<!--fill pages-->
						{{$pageRange := getPageRange .page .showPages .pages}}
						{{range $i, $v := $pageRange}}
							<!--first page-->
							{{if eq $i 0}}
							{{if gt $v 2}}
							<li><a href="/mood">1</a></li>
							<li class="disabled"><a href="javascript:void(0)" style="background-color:#F2F2F2">...</a></li>
							{{end}}
							{{end}}
							
							{{if eq $v $.page}}
							<li class="active"><a href="javascript:void(0);">{{$v}}</a></li>
							{{else}}
								{{if gt $v $.pages}}
									<li class="disabled"><a href="javascript:void(0);">{{$v}}</a></li>
								{{else}}
									<li><a href="/mood?p={{$v}}">{{$v}}</a></li>
								{{end}}
							{{end}}
							
							<!--last page-->
							{{$lastPageIndex := len $pageRange}}
							{{$lastPageIndex := minusInt $lastPageIndex 1}}
							{{$lastPage := minusInt $.pages 1}}
							{{if eq $i $lastPageIndex}}
							{{if lt $v $lastPage}}
							<li class="disabled"><a href="javascript:void(0)"  style="background-color:#F2F2F2">...</a></li>
							<li><a href="/mood?p={{$.pages}}">{{$.pages}}</a></li>
							{{end}}
							{{end}}
						{{end}}
						<!--next page-->
						{{if lt .page .pages}}
						<li><a href="/mood?p={{addInt .page 1}}" aria-lable="Next"><span aria-hidden="true">&raquo;</span></a></li>
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