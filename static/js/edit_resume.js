function showAlert(str) {
	$("#resume-tip").removeClass("hide");
	$("#resume-tip-text").text(str);
}

function submitPostResume(obj) {
	var form = $("#postresume-form");
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