{{define "Title"}}会员 {{.watchedUser.UserName}}{{end}}
{{define "importcss"}}{{end}}
{{define "importjs"}}
<script src="/static/js/project_category.js"></script>
{{end}}
{{define "content"}}
<div id="id-content" class="container">
<!--Modal dialogs-->
	<div id="modalDeleteConfirm" class="modal fade in" role="dialog" aria-hidden="true" style="display: none;">
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<a class="close" data-dismiss="modal">×</a>
					<h4 id="id-modaldeleteconfirm-text"></h4>
				</div>
				<div class="modal-footer">
					<a href="#" onclick="submitDeleteProject(this, '/ajax/project_delete')" class="btn btn-success">确定</a>
					<a href="#" class="btn" data-dismiss="modal">取消</a>
				</div>
			</div>
		</div>
	</div>
	<div id="modalAlert" class="modal fade in" role="dialog" aria-hidden="true" style="display: none;">
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<a class="close" data-dismiss="modal">×</a>
					<h3 style="color:#FE2E2E" id="id-modalalert-text">hint</h3>
				</div>
				<div class="modal-footer">
					<a href="#" class="btn" data-dismiss="modal">关闭</a>
				</div>
			</div>
		</div>
	</div>
	<div id="modalProjectAdd" class="modal fade in" role="dialog" aria-hidden="true" style="display: none;">
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<a class="close" data-dismiss="modal">×</a>
					<h3>添加项目</h3>
				</div>
				<div class="modal-body">
					<form id="id-form-newproject" class="form " novalidate="novalidate" id="new_user" action="/ajax/project_create" accept-charset="UTF-8" method="post">
						<div class="form-group">
						  <input type="input" class="form-control input-lg" placeholder="项目名称" name="project[name]" id="newproject_name" />
						</div>
						<div class="form-group">
						  <input type="input" class="form-control input-lg" placeholder="项目简介" name="project[describe]" id="newproject_describe" />
						</div>
						<div class="form-group">
						  <input type="input" class="form-control input-lg" placeholder="项目封面" name="project[image]" id="newproject_image" />
						</div>
					</form>
				</div>
				<div class="modal-footer">
					<a href="#" onclick="submitCreateProject(this)" class="btn btn-success">添加</a>
					<a href="#" class="btn" data-dismiss="modal">关闭</a>
				</div>
			</div>
		</div>
	</div>
	<div id="modalProjectEdit" class="modal fade in" role="dialog" aria-hidden="true" style="display: none;">
		<div class="modal-dialog">
			<div class="modal-content">
				<div class="modal-header">
					<a class="close" data-dismiss="modal">×</a>
					<h3>编辑项目</h3>
				</div>
				<div class="modal-body">
					<form id="id-form-editproject" class="form " novalidate="novalidate" id="new_user" action="/ajax/project_edit" accept-charset="UTF-8" method="post">
						<div class="form-group">
						  <input type="input" class="form-control input-lg" placeholder="项目名称" name="project[name]" id="editproject_name" />
						</div>
						<div class="form-group">
						  <input type="input" class="form-control input-lg" placeholder="项目简介" name="project[describe]" id="editproject_describe" />
						</div>
						<div class="form-group">
						  <input type="input" class="form-control input-lg" placeholder="项目封面" name="project[image]" id="editproject_image" />
						</div>
					</form>
				</div>
				<div class="modal-footer">
					<a href="javascript:void(0);" onclick="submitEditProject(this)" class="btn btn-success">编辑</a>
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
			</p>
			{{end}}
			{{range .category}}
			<div id="id-div-{{.ProjectName}}" class="panel panel-default">
				<div class="panel-heading" style="text-align:center;">
					<a href="/project/{{.ProjectName}}/1">{{.ProjectName}}</a>
					{{if gt $.user.Permission 3}}
					<div style="float:right;margin-bottom:2px">
						<button id="id-project-modify-{{.ProjectName}}" onclick="onEditProject(this)" type="button" project="{{.ProjectName}}" class="btn btn-sm btn-primary">编辑项目</button>
						<button id="id-project-del-{{.ProjectName}}" onclick="onDelProject(this)" type="button" project="{{.ProjectName}}" class="btn btn-sm btn-danger">删除项目</button>
					</div>
					{{end}}
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