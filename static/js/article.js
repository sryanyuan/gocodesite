$(document).ready(function(){
	
})

function showAlert(text) {
	$("#id-modalalert-text").html(text);
	$("#modalAlert").modal({backdrop:"static"});
}

function deleteArticle(obj, articleId) {
	var articleTitle = $(obj).attr("articleTitle");
	
	//	show modal
	$("#id-modaldeleteconfirm-text").html("您确认要删除文章["+articleTitle+"] ?");
	$("#id-modaldeleteconfirm-text").attr("articleId", articleId);
	$("#modalDeleteConfirm").modal({backdrop:"static"});
}

function submitDeleteArticle(obj, action) {
	var articleId = $("#id-modaldeleteconfirm-text").attr("articleId");
	var postData = 'articleId='+articleId;

	$.post(action, postData, function(ret){
		if (ret.Result == 0) {
			location.href = ret.Msg;
		} else {
			$("#modalDeleteConfirm").modal("hide");
			showAlert(ret.Msg);
		}
	}).error(function(e){
		$("#modalDeleteConfirm").modal("hide");
		showAlert("请求超时");
	});
}

function topArticle(obj, top, articleId) {
	var action = "/ajax/article_top";
	var postData = "top=1&articleId="+articleId;
	
	if (!top) {
		postData = "top=0&articleId="+articleId;
	}
	$.post(action, postData, function(ret){
		if (ret.Result == 0) {
			location.reload();
		} else {
			showAlert(ret.Msg);
		}
	}).error(function(e){
		showAlert("请求超时");
	});
}