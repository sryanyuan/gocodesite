function submitReply(sender) {
    var form = $("#reply_form");

    $.post(form.attr("action"), form.serialize(), function(ret){
		if (ret.Result == 0) {
			//	ok, refresh
			location.href = ret.Msg;
		} else {
            // Refresh captcha id
            if(null != ret.CaptchaId &&
				ret.CaptchaId.length != 0){
				var signInCaptcha = $("#id-signup-captchaimg")
				signInCaptcha.attr("src", "/captcha/"+ret.CaptchaId+".png");
				$("#id-signup-captchaIdHolder").attr("value", ret.CaptchaId)
			}
			//	failed
			alert(ret.Msg);
		}
	}).error(function(e){
		alert(e.responseText);
	});
}

$(document).ready(function(){
    var signInCaptcha = $("#id-signup-captchaimg")
	signInCaptcha.click(function(event){
		event.preventDefault();
		var refurl = $(this).attr("src")
		$(this).attr("src", refurl+"?reload="+(new Date()).getTime());
	});
})

function deleteReply(obj, articleId) {
    var replyId = $(obj).attr("replyId");
	// Show modal
	$("#id-modalreplydeleteconfirm-text").html("您确认要删除该条回复?");
	$("#id-modalreplydeleteconfirm-text").attr("replyId", replyId);
	$("#modalReplyDeleteConfirm").modal({backdrop:"static"});
}

function submitDeleteReply(obj, action) {
	var replyId = $("#id-modalreplydeleteconfirm-text").attr("replyId");
	var postData = 'replyId='+replyId;

	$.post(action, postData, function(ret){
		if (ret.Result == 0) {
            $("#modalReplyDeleteConfirm").modal("hide");
            location.reload();
			//location.href = ret.Msg;
		} else {
			$("#modalReplyDeleteConfirm").modal("hide");
			alert(ret.Msg);
		}
	}).error(function(e){
		$("#modalReplyDeleteConfirm").modal("hide");
		alert("请求超时");
	});
}