$("#id-project-add").click(function(){
	$("#modalProjectAdd").modal({backdrop:"static"});
})

$("#id-project-modify").click(function(){
	$("#modalProjectEdit").modal({backdrop:"static"});
	
	//	fill content
})

function showAlert(text) {
	$("#id-modalalert-text").html(text);
	$("#modalAlert").modal({backdrop:"static"});
}

function onEditProject(obj, projectName) {
	var project = $(obj).attr("project");
	var projectId = $(obj).attr("projectId");
	//	get content
	var desc = $("#id-div-"+projectId).find("p").html();
	$("#input-projectid").val(projectId);
	//	get image
	var image = $(obj).attr("projectImage");
	
	//	show modal
	$("#editproject_name").val(project);
	$("#editproject_describe").val(desc);
	$("#editproject_image").val(image);
	$("#modalProjectEdit").modal({backdrop:"static"});
}

function onDelProject(obj) {
	var project = $(obj).attr("project");
	var projectId = $(obj).attr("projectId");
	
	//	show modal
	$("#id-modaldeleteconfirm-text").html("您确认要删除项目["+project+"] ?");
	$("#id-modaldeleteconfirm-text").attr("projectId", projectId);
	$("#modalDeleteConfirm").modal({backdrop:"static"});
}

$(document).ready(function(){
	
})

function submitCreateProject(obj) {
	var form = $("#id-form-newproject")
	$.post(form.attr("action"), form.serialize(), function(ret){
		if (ret.Result == 0) {
			//	ok, refresh
			location.reload();
		} else {
			//	failed
			$("#modalProjectAdd").modal("hide");
			showAlert(ret.Msg);
		}
	}).error(function(e){
		$("#modalProjectAdd").modal("hide");
		showAlert("请求超时");
	});
}

function submitDeleteProject(obj, action) {
	var projectId = $("#id-modaldeleteconfirm-text").attr("projectId");
	var postData = 'project[id]='+projectId;

	$.post(action, postData, function(ret){
		if (ret.Result == 0) {
			location.reload();
		} else {
			$("#modalDeleteConfirm").modal("hide");
			showAlert(ret.Msg);
		}
	}).error(function(e){
		$("#modalDeleteConfirm").modal("hide");
		showAlert("请求超时");
	});
}

function submitEditProject(obj) {
	var form = $("#id-form-editproject")
	$.post(form.attr("action"), form.serialize(), function(ret){
		if (ret.Result == 0) {
			//	ok, refresh
			location.reload();
		} else {
			//	failed
			$("#modalProjectEdit").modal("hide");
			showAlert(ret.Msg);
		}
	}).error(function(e){
		$("#modalProjectEdit").modal("hide");
		showAlert("请求超时");
	});
}