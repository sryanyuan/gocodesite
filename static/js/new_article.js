$(document).ready(function(){
	var signInCaptcha = $("#id-article-captchaimg")
	signInCaptcha.click(function(event){
		event.preventDefault();
		var refurl = $(this).attr("src")
		$(this).attr("src", refurl+"?reload="+(new Date()).getTime());
	});
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
			
			//	refresh captcha?
			if(null != ret.CaptchaId &&
				ret.CaptchaId.length != 0){
				var signInCaptcha = $("#id-article-captchaimg")
				signInCaptcha.attr("src", "/captcha/"+ret.CaptchaId+".png");
				$("#id-article-captchaIdHolder").attr("value", ret.CaptchaId)
			}
		}
	}).error(function(e){
		showAlert("请求超时");
	});
}