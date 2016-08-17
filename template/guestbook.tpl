{{define "Title"}}关于{{end}}
{{define "importcss"}}
<link href="/static/css/guestbook.css" rel="stylesheet" />
{{end}}
{{define "importjs"}}
<script src="/static/js/comment_load.js"></script>
<script src="/static/js/comment_recent_visitors.js"></script>
{{end}}
{{define "content"}}
<div id="id-content" class="container">
	<div class="row">
		<div class="col-md-8 col-md-offset-0">
			<ul class="breadcrumb">
				<li><a href="/"><i class="fa fa-home"></i> 首页</a></li>
				<li class="active">留言板</li>
			</ul>
			<div id="comment-container" class="shadow-box white-box">
				{{template "comment_guestbook_html" .}}
			</div>
		</div>
		<div class="col-md-3 col-md-offset-0">
			<h2 class="section-title-s2">最近访客</h2>
			<div id="comment-visitor-container" class="white-box shadow-box">
				<ul class="ds-recent-visitors" data-num-items="30" id="ds-recent-visitors"></ul>
			</div>
		</div>
	</div>
</div>
{{end}}