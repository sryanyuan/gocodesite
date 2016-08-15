function showAlert(text) {
	$("#id-modalalert-text").html(text);
	$("#modalAlert").modal({backdrop:"static"});
}

function uploadFile(obj) {
	var form = $("#upload-form")
	$("#id-modalalert-text").css("color", "#FE2E2E");
	var inputFile = $("#selected-file").val();
	if (inputFile.length == 0) {
		showAlert("请选择文件");
		return;
	}
	
	var ret = $("#submit-upload").click();
	
	/*$.post(form.attr("action"), form.serialize(), function(ret){
		if (ret.Result == 0) {
			//	ok, refresh
			$("#id-modalalert-text").css("color", "#ffffff");
			showAlert("上传成功");
		} else {
			//	failed
			showAlert(ret.Msg);
		}
	}).error(function(e){
		showAlert("请求超时");
	});*/
}