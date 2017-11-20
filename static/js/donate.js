function showMsgTip(id, err, msg) {
	var widget = $(id);
	if (nil == widget) {
		console.debug("id ", id, " not found");
		return
	}
	if (err == 0) {
		widget.removeClass("alert-danger");
		widget.addClass("alert-info");
	} else {
		widget.removeClass("alert-info");
		widget.addClass("alert-danger");
	}
	widget.removeClass("hidden");
	widget.html(ret.Msg);
}

function changeAlertLook(widget, err) {
	widget.removeClass("alert-danger");
	widget.removeClass("alert-info");
	if (err == 0) {
		$("#charge-hint-title").html("");
		widget.addClass("alert-info");
	} else {
		$("#charge-hint-title").html("错误");
		widget.addClass("alert-danger");
	}
}

$(document).ready(function(){
    var donateForm = $("#id-form-charge");
    var donateBtnID = "#id-charge-zfbqr";
	if(null != donateForm){
		donateForm.submit(function(event){
			event.preventDefault();
			var target = event.target;
			var action = $(target).attr("action");
			if($(donateBtnID).hasClass("disabled")){
				return;
			}

			$(donateBtnID).addClass("disabled");
			$.post(action, $(target).serialize(), function(ret){
				$("#id-charge-submit").removeClass("disabled");
				var chargeHint = $("#id-charge-hint")
				if(null != ret){
					if(0 != ret.Result){
						chargeHint.removeClass("hidden");
						changeAlertLook(chargeHint, 1);
                        $("#id-charge-hinttext").html(ret.Msg);
                        $(donateBtnID).removeClass("disabled");
					}else{
						chargeHint.removeClass("hidden");
						changeAlertLook(chargeHint, 0);
                        $("#id-charge-hinttext").html("订单生成成功，请用支付宝钱包扫码支付，成功后请不要关闭本页面，直到通知");
						//location.href = ret.Msg;
					}
				}
			}).error(function(e){
				$("#id-charge-submit").removeClass("disabled");
				$("#id-charge-hint").removeClass("hidden");
				changeAlertLook(chargeHint, 1);
				$("#id-charge-hinttext").html("请求失败，请检查网络");
			});
		})
	}
});