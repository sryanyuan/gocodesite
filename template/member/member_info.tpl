{{define "Title"}}会员 {{.watchedUser.NickName}}{{end}}
{{define "importcss"}}
<link href="/static/css/member_info.css" rel="stylesheet" />
{{end}}
{{define "importjs"}}
<script src="/static/js/member_info.js"></script>
{{end}}
{{define "content"}}
<div id="id-content" class="container">
	<div class="row">
		<div class="sidebar col-md-4">
			<div class="panel panel-default">
			  <div class="panel-body">
				<div class="media">
					<div class="media-left">
						<img alt="{{.user.UserName}}" class="media-object avatar-120 img-circle" width="72" height="72" src="{{.imgPrefix}}/male.png" />
					</div>
					<div id="member-baseinfo" class="media-body">
						<div class="item username">
							{{.watchedUser.UserName}}(<span id="id-meminfo-nickname">{{.watchedUser.NickName}})
						</div>
						<div class="item" id="id-meminfo-signupseq">
							第 {{.watchedUser.Uid}} 位会员 / {{getUnixTimeString .watchedUser.CreateTime}}
						</div>
						<div class="item" id="id-meminfo-postinfo">
							发了 {{.postCount}} 贴 • 回了 {{.replyCount}} 贴
						</div>
						<div class="item social">
							{{if ne .watchedSocialInfo.Weibo ""}}
							<a href="{{.watchedSocialInfo.Weibo}}"><i class="fa fa-weibo"></i></a>
							{{end}}
							{{if ne .watchedSocialInfo.Github ""}}
							<a href="{{.watchedSocialInfo.Github}}"><i class="fa fa-github"></i></a>
							{{end}}
						</div>
					</div>
					<div id="id-member-following">
						<div id="member-mood">
							{{$moodLen := len .watchedUser.Mood}}
							{{if ne $moodLen 0}}
							{{.watchedUser.Mood}}
							{{else}}
							该用户很懒，没有设置签名
							{{end}}
						</div>
					</div>
				</div>
			  </div>
			</div>
			{{if eq .watchedUser.Uid .user.Uid}}
			{{if gt .user.Permission 3}}
			<div class="panel panel-default">
				<div class="panel-heading" style="text-align:center;">
					<h3 class="panel-title">管理面板</h3>
				</div>
				<div class="panel-body">
					<div>
						<a href="/admin/upload" class="btn btn-sm btn-success">上传静态文件</a>
					</div>
				</div>
			</div>
			{{end}}
			{{end}}
		</div>
		<div class="col-md-8">
			<!--div class="panel panel-default">
				<div id="id-member-intro" class="panel-heading">
					个人介绍
				</div>
				<div class="panel-body">
					<div>
						
					</div>
				</div>
			</div!-->
			<div>
				<ul class="nav nav-tabs">
					<li class="active"><a href="#id-member-tab-post" data-toggle="tab">最近主题</a></li>
					<li><a href="#id-member-tab-reply" data-toggle="tab">最近回复</a></li>
				</ul>
				<div class="tab-content">
					<div class="tab-pane active post fade in" id="id-member-tab-post">
						<div class="panel panel-default">
							<ul id="member-post-list-group" class="list-group" articleCount="{{len .articles}}">
								{{range $i, $v := .articles}}
								<li class="list-group-item">
									<div class="recent-article-title">
										<a href="/project/{{$v.ProjectId}}/article/{{$v.Id}}">{{$v.ArticleTitle}}</a>
									</div>
									<div class="recent-article-info">
										<span class="member-post-info">发表于 {{getTimeGapString $v.PostTime}}</span>
										   • <span class="member-post-reply-count" id="member-post-reply-count-{{$v.Id}}" articleId="{{$v.Id}}">0</span> 个回复
									</div>
								</li>
								{{end}}
								<div style="text-align:center;margin-top:10px;padding-bottom:10px;">
									<button type="button" class="btn btn-link" style="background-color:white; border:1px solid #DDDDDD;border-radius:4px;"><a href="/member/{{.watchedUser.NickName}}/articles?p=1">查看所有主题</a></button>
								</div>
							</ul>
						</div>
					</div>
					<div class="tab-pane fade" id="id-member-tab-reply">
						<p>Nothing</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
{{end}}