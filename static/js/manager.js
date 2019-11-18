function showAlert(text) {
	$("#id-modalalert-text").html(text);
	$("#modalAlert").modal({backdrop:"static"});
}

$("#bmkv-add").click(function(){
	$("#modalKvAdd").modal({backdrop:"static"});
})

function bmkvupdate(obj){
    $("#bmkvupdate-key").val($(obj).parent().siblings("#kv-key").html());
    $("#bmkvupdate-value").val($(obj).parent().siblings("#kv-value").html());
	$("#modalKvUpdate").modal({backdrop:"static"});
}

function bmkvdelete(obj){
    $("#bmkvdel-key").val($(obj).parent().siblings("#kv-key").html());
	$("#modalKvDelete").modal({backdrop:"static"});
}

function submitBmkvAdd(form) {
    form = $("#id-form-kvadd")
    $.post(form.attr("action"), form.serialize(), function(ret){
		if (ret.Result == 0) {
			//	ok, refresh
			location.reload();
		} else {
			//	failed
			$("#modalKvAdd").modal("hide");
			showAlert(ret.Msg);
		}
	}).error(function(e){
		$("#modalKvAdd").modal("hide");
		showAlert("请求超时");
	});
}

function submitBmkvUpdate(form) {
    form = $("#id-form-kvupdate")
    $.post(form.attr("action"), form.serialize(), function(ret){
		if (ret.Result == 0) {
			//	ok, refresh
			location.reload();
		} else {
			//	failed
			$("#modalKvAdd").modal("hide");
			showAlert(ret.Msg);
		}
	}).error(function(e){
		$("#modalKvAdd").modal("hide");
		showAlert("请求超时");
	});
}

function submitBmkvDelete(form) {
    form = $("#id-form-kvdelete")
    $.post(form.attr("action"), form.serialize(), function(ret){
		if (ret.Result == 0) {
			//	ok, refresh
			location.reload();
		} else {
			//	failed
			$("#modalKvDelete").modal("hide");
			showAlert(ret.Msg);
		}
	}).error(function(e){
		$("#modalKvDelete").modal("hide");
		showAlert("请求超时");
	});
}