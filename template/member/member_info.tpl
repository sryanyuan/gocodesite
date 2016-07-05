{{define "Title"}}会员 {{.watchedUser.UserName}}{{end}}
{{define "importcss"}}{{end}}
{{define "importjs"}}{{end}}
{{define "content"}}
<div id="id-content" class="container">
	<div class="row">
	<div class="sidebar col-md-4">
		<div>
		<div class="panel">
          <div class="panel panel-heading">
			{{.watchedUser.UserName}}
			{{if .isSelf}}
			<span id="id-member-signupseq">您是本站的第{{.watchedUser.Uid}}位会员</span>
			{{else}}
			<span id="id-member-signupseq">本站的第{{.watchedUser.Uid}}位会员</span>
			{{end}}
			<div class="panel-body">
				<img alt="{{.user.UserName}}" class="avatar img-rounded" height="200" src="{{.imgPrefix}}/male.png" width="200" />
			</div>
		  </div>
		</div>
		</div>
	</div>
	</div>
</div>
{{end}}