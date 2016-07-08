{{define "Title"}}会员 {{.watchedUser.UserName}}{{end}}
{{define "importcss"}}{{end}}
{{define "importjs"}}{{end}}
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
					<div class="media-body">
						<div class="item username">
							{{.watchedUser.UserName}}(<span id="id-meminfo-nickname">{{.watchedUser.NickName}})
						</div>
						<div class="item" id="id-meminfo-signupseq">
							第 {{.watchedUser.Uid}} 位会员 / {{getUnixTimeString .watchedUser.CreateTime}}
						</div>
						<div class="item" id="id-meminfo-postinfo">
							发了 {{.postCount}} 贴  回了 {{.replyCount}} 贴
						</div>
						<div class="item social">
						</div>
					</div>
				</div>
			  </div>
			</div>
		</div>
	</div>
</div>
{{end}}