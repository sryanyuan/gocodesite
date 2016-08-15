{{define "Title"}}上传静态文件{{end}}
{{define "importcss"}}{{end}}
{{define "importjs"}}{{end}}
{{define "content"}}
<div id="id-content" class="container">
	<div class="row">
		<div class="col-md-6 col-md-offset-3">
			<div style="text-align:center;">
				{{if eq .Result ""}}
				<h1 style="color:green;">
				{{else}}
				<h1 style="color:red;">
				{{end}}
				{{.Title}}
				</h1>
				<span style="font:16px;">{{.Text}}</span>
				<hr/>
				<div>
					<a href="javascript :;" onclick="javascript :history.back(-1);" class="btn btn-success">回到上一页</a>
					<a style="margin-left:10px;" href="/" class="btn btn-success">回到主页</a>
				</div>
			</div>
		</div>
	</div>
</div>
{{end}}