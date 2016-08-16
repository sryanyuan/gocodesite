{{define "Title"}}分类{{end}}
{{define "importcss"}}
<link href="/static/css/project-category.css" rel="stylesheet" />
{{end}}
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
						  <input type="input" class="form-control input-md" placeholder="项目名称" name="project[name]" id="editproject_name" />
						</div>
						<div class="form-group">
						  <input type="input" class="form-control input-md" placeholder="项目简介" name="project[describe]" id="editproject_describe" />
						</div>
						<div class="form-group">
						  <input type="input" class="form-control input-md" placeholder="项目封面" name="project[image]" id="editproject_image" />
						</div>
						<input id="input-projectid" type="hidden" name="project[id]" value="" />
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
		<div class="col-md-8 col-md-offset-2">
			<!--Administrator panel-->
			{{if gt .user.Permission 3}}
			<p>
				<button id="id-project-add" type="button" class="btn btn-sm btn-primary">添加项目</button>
			</p>
			{{end}}
			
			<div id="category-container">
				{{range .category}}
				<div class="media category-box">
					<a class="pull-left" href="/project/{{.Id}}/page/1" target="_blank">
						<img class="media-object" src="
						{{if eq .Image ""}}
						{{$.imgPrefix}}/category_cover.png
						{{else}}
						{{$.imgPrefix}}/{{.Image}}
						{{end}}"
						width="215" height="144" />
					</a>
					<div class="media-body">
						<h3>
							<a href="/project/{{.Id}}/page/1" target="_blank">{{.ProjectName}}</a>
						</h3>
						<div class="category-info">
							<i class="fa fa-smile-o"></i> 创建者：<a href="/member/{{.Author}}" target="_blank">{{.Author}}</a>
							<span style="margin-left:8px;"><i class="fa fa-book"></i> 文章数：{{.ItemCount}}<span>
							<span style="margin-left:8px;"><i class="fa fa-universal-access"></i> 发帖权限：
							{{if lt .PostPriv 4}}
								普通用户
							{{else}}
								超级管理员
							{{end}}
							<span>
						</div>
						<div style="border-bottom:1px solid #EEEEEE;margin-top:8px;margin-bottom:5px;"></div>
						<div id="category-desc-{{.Id}}" class="category-desc">
							{{.ProjectDescribe}}
						</div>
					</div>
					{{if gt $.user.Permission 3}}
					<div style="float:right;margin-bottom:2px">
						<button id="id-project-modify-{{.Id}}" onclick="onEditProject(this,{{.Id}})" type="button" projectId="{{.Id}}" projectImage="{{.Image}}" project="{{.ProjectName}}" class="btn btn-sm btn-primary">编辑项目</button>
						<button id="id-project-del-{{.Id}}" onclick="onDelProject(this)" type="button" project="{{.ProjectName}}" projectId="{{.Id}}" class="btn btn-sm btn-danger">删除项目</button>
					</div>
					{{end}}
				</div>
				{{end}}
			</div>
		</div>
	</div>
</div>
{{end}}