$(document).ready(function(){
	
})

function showAlert(str) {
	$("#article-tip").removeClass("hide");
	$("#article-tip-text").text(str);
}

function submitPostArticle(obj) {
	var form = $("#postarticle-form");
	if ($("#text-title").val().length == 0) {
		showAlert("请输入标题");
		return;
	}
	if (editor.getMarkdown().trim().length == 0) {
		showAlert("请输入内容");
		return;
	}
	
	$.post(form.attr("action"), form.serialize(), function(ret){
		if (ret.Result == 0) {
			//	ok, refresh
			location.href = ret.Msg;
		} else {
			//	failed
			showAlert(ret.Msg);
		}
	}).error(function(e){
		showAlert("请求超时");
	});
}