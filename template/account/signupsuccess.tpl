{{define "Title"}}关于{{end}}
{{define "importcss"}}{{end}}
{{define "importjs"}}{{end}}
{{define "content"}}
<div id="id-content" class="row">
	<div class="col-md-8 col-md-offset-2">
		<div id="id-signup-hint" class="alert alert-success" role="alert">
          <strong>恭喜，您的账户({{.account}})已注册成功，请牢记密码。请<a href="/account/signin">登陆</a>
		</div>
	</div>
</div>
{{end}}