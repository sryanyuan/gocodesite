//	菜单
$('li.dropdown').mouseover(function() {   
     $(this).addClass('open');    }).mouseout(function() {        $(this).removeClass('open');    }); 
	 
//	提交注册信息
var signUpForm = $("#id-form-signup");
if(null != signUpForm)
{
	signUpForm.submit(function(event){
		event.preventDefault();
		var target = event.target;
		var action = $(target).attr("action");
		if($("#id-signup-submit").hasClass("disabled"))
		{
			return;
		}
		//	set disable mode
		$("#id-signup-submit").addClass("disabled");
		$.post(action, $(target).serialize(), function(ret){
			$("#id-signup-submit").removeClass("disabled");
			var signUpHint = $("#id-signup-hint")
			if(null != ret.Result){
				if(0 != ret.Result){
					signUpHint.removeClass("hidden");
					$("#id-signup-hinttext").html(ret.Msg);
				}else{
					location.href = ret.Msg;
				}
			}
		}).error(function(e){
			$("#id-signup-submit").removeClass("disabled");
			$("#id-signup-hint").removeClass("hidden");
			$("#id-signup-hinttext").html("请求失败，请检查网络");
		});
	})
}

//	提交登陆信息
var signInForm = $("#id-form-signin");
if(null != signInForm){
	signInForm.submit(function(event){
		event.preventDefault();
		var target = event.target;
		var action = $(target).attr("action");
		if($("#id-signin-submit").hasClass("disabled")){
			return;
		}

		$("#id-signin-submit").addClass("disabled");
		$.post(action, $(target).serialize(), function(ret){
			$("#id-signin-submit").removeClass("disabled");
			var signInHint = $("#id-signin-hint")
			if(null != ret){
				if(0 != ret.Result){
					signInHint.removeClass("hidden");
					$("#id-signin-hinttext").html(ret.Msg);
				}else{
					location.href = ret.Msg;
				}
			}
		}).error(function(e){
			$("#id-signin-submit").removeClass("disabled");
			$("#id-signin-hint").removeClass("hidden");
			$("#id-signin-hinttext").html("请求失败，请检查网络");
		});
	})
}else{
	alert("nil form");
}