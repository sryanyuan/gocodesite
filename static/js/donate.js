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

function show_pay_tip() {
	$("#id-modalalert-text").html("特别注意：请输入系统实际需要支付的金额，否则点数将无法到账。比如您捐助10点，但是系统显示的支付金额为9.99，则请支付9.99，不要多付，成功后您的到账将依旧是10点");
	$("#modalAlert").modal({backdrop:"static"});
}

$("#id-charge-zfbqr").click(function zfbpay(event){
	event.preventDefault();
	$("#id-pay-method").val("0");
	show_pay_tip();
});

$("#id-charge-wxqr").click(function wxpay(event) {
	event.preventDefault();
	$("#id-pay-method").val("1");
	show_pay_tip();
});

function pay_next() {
	var donateForm = $("#id-form-charge");
	donateForm.submit();
}

function unionpay() {
	$("#id-pay-method").val("2");
	var donateForm = $("#id-form-charge");
	donateForm.submit();
}

var timeHandle = null;

$(document).ready(function(){
    var donateForm = $("#id-form-charge");
	var donateBtnID = "#id-charge-zfbqr";
	var donateWxBtnID = "#id-charge-wxqr";
	if(null != donateForm){
		donateForm.submit(function(event){
			event.preventDefault();
			var target = event.target;
			var action = $(target).attr("action");
			if($(donateBtnID).hasClass("disabled")){
				return;
			}
			if($(donateWxBtnID).hasClass("disabled")){
				return;
			}

			$(donateBtnID).addClass("disabled");
			$(donateWxBtnID).addClass("disabled");
			$.post(action, $(target).serialize(), function(ret){
				$("#id-charge-submit").removeClass("disabled");
				var chargeHint = $("#id-charge-hint");
				if(null != ret){
					if(0 != ret.Result){
						chargeHint.removeClass("hidden");
						changeAlertLook(chargeHint, 1);
                        $("#id-charge-hinttext").html(ret.Msg);
						$(donateBtnID).removeClass("disabled");
						$(donateWxBtnID).removeClass("disabled");
					}else{
						chargeHint.removeClass("hidden");
						changeAlertLook(chargeHint, 0);

						var orderInfo = JSON.parse(ret.Msg);
						// Show pay iframe
						var payMethod = $("#id-pay-method").val();
						var payWindow = $("#alipay_qr_iframe");
						var paysrc = "";
						if ("0" == payMethod || "1" == payMethod) {
							var payName = "支付宝钱包";
							if ("1" == payMethod) {
								payName = "微信";
							}
							paysrc = orderInfo.PpayURL + "/static/payPage/pay.html?orderId=" + orderInfo.PpayOrderID;
							$("#id-charge-hinttext").html("订单号[" + orderInfo.PpayOrderID + "] (请牢记)，请用" + payName + "扫码支付，成功后请不要关闭本页面，直到跳转");
						} else if ("2" == payMethod) {
							$("#id-charge-hinttext").html("订单号[" + orderInfo.OrderID + "] (请牢记)，请扫码支付，成功后请不要关闭本页面，直到跳转");
						} else {
							changeAlertLook(chargeHint, 1);
							$("#id-charge-hinttext").html("非法的url");
							return;
						}

						if ("1" == payMethod || "0" == payMethod) {
							payWindow.attr("src", paysrc);
							payWindow.removeClass("hidden");
							window.location.href = paysrc;
						} else if ("2" == payMethod) {
							// Draw QR code
							$("#pay_qrcode").qrcode(orderInfo.QRUrl);
						}

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
	var dt = new Date();
	var rurl = "/donate/" + orderID + "?ts=" + dt.getTime();
	var chargeHint = $("#id-charge-hint");

	$.get(rurl, function(result) {
		var orderStatus = JSON.parse(result);
		if (orderStatus.Msg == "OK") {
			changeAlertLook(chargeHint, 0);
			$("#id-charge-hinttext").html("订单支付成功 [" + orderID + "]");
			clearInterval(timeHandle);
			$(donateWxBtnID).removeClass("disabled");
			$(donateBtnID).removeClass("disabled");
		} else if (orderStatus.Msg.indexOf("wait")) {
			// Nothing
		} else {
			changeAlertLook(chargeHint, 1);
			$("#id-charge-hinttext").html("订单支付失败 [" + orderID + "] " + orderStatus.Msg + " , 请联系管理员");
			clearInterval(timeHandle);
		}
	})
}
