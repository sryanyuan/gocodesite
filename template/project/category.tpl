{{define "Title"}}会员 {{.watchedUser.UserName}}{{end}}
{{define "importcss"}}{{end}}
{{define "importjs"}}
<script src="/static/js/project_category.js"></script>
{{end}}
{{define "content"}}
<div id="id-content" class="container">
<!--Modal dialogs-->
	<div id="modalProjectAdd" class="modal fade in" role="dialog" aria-hidden="true" style="display: none;">
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<a class="close" data-dismiss="modal">×</a>
					<h3>添加项目</h3>
				</div>
				<div class="modal-body">
					<form id="id-form-newproject" class="form " novalidate="novalidate" id="new_user" action="/manager/newproject" accept-charset="UTF-8" method="post">
						<div class="form-group">
						  <input type="email" class="form-control input-lg" placeholder="项目名称" name="project[name]" id="project_name" />
						</div>
						<div class="form-group">
						  <input type="password" class="form-control input-lg" placeholder="项目简介" name="project[describe]" id="project_describe" />
						</div>
						<div class="form-group">
						  <input type="password" class="form-control input-lg" placeholder="项目封面" name="project[image]" id="project_image" />
						</div>
					</form>
				</div>
				<div class="modal-footer">
					<a href="#" class="btn btn-success">添加</a>
					<a href="#" class="btn" data-dismiss="modal">关闭</a>
				</div>
			</div>
		</div>
	</div>
	<div class="row">
		<div class="col-md-6 col-md-offset-3">
			{{if gt .user.Permission 3}}
			<!--Administrator panel-->
			<p>
				<button id="id-project-add" type="button" class="btn btn-sm btn-primary">添加项目</button>
				<button id="id-project-remove" type="button" class="btn btn-sm btn-primary">删除项目</button>
				<button id="id-project-modify" type="button" class="btn btn-sm btn-primary">编辑项目</button>
			</p>
			{{end}}
			{{range .category}}
			<div class="panel panel-default">
				<div class="panel-heading" style="text-align:center;">
					<a href="/project/{{.ProjectName}}">{{.ProjectName}}</a>
				</div>
			  <div class="panel-body">
				<p>{{.ProjectDescribe}}</p>
			  </div>
			</div>
			{{end}}
		</div>
	</div>
</div>
{{end}}