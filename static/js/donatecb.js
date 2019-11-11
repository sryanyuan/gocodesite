function getQueryVariable(variable)
{
       var query = window.location.search.substring(1);
       var vars = query.split("&");
       for (var i=0;i<vars.length;i++) {
               var pair = vars[i].split("=");
               if(pair[0] == variable){return pair[1];}
       }
       return(false);
}

$("#id-charge-hint").removeClass("alert-danger");
$("#id-charge-hint").removeClass("alert-info");
$("#id-charge-hint").removeClass("hidden");
$("#charge-hint-title").html("");
$("#id-charge-hint").addClass("alert-info");

$("#id-charge-hinttext").html("订单支付成功 [" + getQueryVariable("payId") + "]，实际支付 " + 
getQueryVariable("reallyPrice") + "，实际到账 " + getQueryVariable("price") + " 点");