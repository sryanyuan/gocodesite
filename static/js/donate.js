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

var timeHandle = null;

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
				var chargeHint = $("#id-charge-hint");
				if(null != ret){
					if(0 != ret.Result){
						chargeHint.removeClass("hidden");
						changeAlertLook(chargeHint, 1);
                        $("#id-charge-hinttext").html(ret.Msg);
                        $(donateBtnID).removeClass("disabled");
					}else{
						chargeHint.removeClass("hidden");
						changeAlertLook(chargeHint, 0);

						var orderInfo = JSON.parse(ret.Msg);
                        $("#id-charge-hinttext").html("订单号<" + orderInfo.OrderID + "> (请牢记)，请用支付宝钱包扫码支付，成功后请不要关闭本页面，直到跳转");
						//location.href = ret.Msg;
						// Show pay iframe
						var payWindow = $("#alipay_qr_iframe");
						var paysrc = "https://api.jsjapp.com/plugin.php?id=add:alipay2&addnum=" + orderInfo.OrderID + "&total=" + orderInfo.NumFloat + "&apiid=" + orderInfo.ApiID + "&apikey=" + orderInfo.ApiKey + "&uid=" + orderInfo.Uid + "&showurl=";
						var cburl = orderInfo.CallHost + "/ctrl?cmd=insertdonatecb&secret=" + orderInfo.CallSecret;
						paysrc = paysrc + cburl;
						payWindow.attr("src", paysrc);
						payWindow.removeClass("hidden");

						// Payment result check
						timeHandle = window.setInterval(wrapPaymentResult(null, orderInfo.OrderID, orderInfo.CallHost), 2000);
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

function wrapPaymentResult(timer, orderID, calladdr) {
	return function() {
		_checkPaymentResult(timer, orderID, calladdr);
	}
}

function _checkPaymentResult(timer, orderID, calladdr) {
	var rurl = "/donate/" + orderID;
	var chargeHint = $("#id-charge-hint");

	$.get(rurl, function(result) {
		var orderStatus = JSON.parse(result);
		if (orderStatus.Msg == "OK") {
			changeAlertLook(chargeHint, 0);
			$("#id-charge-hinttext").html("订单支付成功 <" + orderID + ">");
			clearInterval(timeHandle);
		} else if (orderStatus.Msg.indexOf("wait")) {
			// Nothing
		} else {
			changeAlertLook(chargeHint, 1);
			$("#id-charge-hinttext").html("订单支付失败 <" + orderID + "> " + orderStatus.Msg + " , 请联系管理员");
			clearInterval(timeHandle);
		}
	})
}
