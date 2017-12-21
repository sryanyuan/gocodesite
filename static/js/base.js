//	菜单
$('li.dropdown').mouseover(function() {   
     $(this).addClass('open');    }).mouseout(function() {        $(this).removeClass('open');    }); 

function adjustFooter(){
	return;
	
	if ($(window).height() != $(document).height()) {
		$("#id_footer").removeClass("navbar-fixed-bottom");
	} else {
		if (!$("#id_footer").hasClass("navbar-fixed-bottom")) {
			$("#id_footer").addClass("navbar-fixed-bottom");
		}
	}
}

function adjustBodyMinHeight() {
	if($(window).height() == $(document).height()){
		var height = $(window).height() - $("#id_footer").height();
		$("body").css("min-height", height+"px");
	} else {
		$("body").css("min-height", "none");
	}
}

var totalMessageCnt = 0;

function pullMessageCount() {
	$.getJSON("/ajax/message_get_count", function(ret){
		if (0 == ret.Result) {
			if (null == ret.Msg || 
			ret.Msg.length == 0 ) {
				return;
			}

			var cnt = parseInt(ret.Msg);
			if (cnt == 0 || isNaN(cnt)) {
				return;
			}
			totalMessageCnt = cnt;
			// Add tip to navbar
			$("#navbar_message").removeClass("hidden");
			$("#navbar_message_count").html(ret.Msg);
		}
	})
}

function formatMessageHTML() {
	// Pull all message content
	$.getJSON("/ajax/message_get", function(ret){
		if (0 == ret.Result) {
			if (null == ret.Msg || 
			ret.Msg.length == 0 ) {
				return;
			}

			var messages = JSON.parse(ret.Msg);
			if (messages.length == 0) {
				return;
			}
			// Add content
			setTimeout(function(){
				var container = $("#id_message_pop_container");
				for (var i in messages) {
					var action = ""
					if (messages[i].Type == 1) {
						action = " 评论了 "
					} else if (messages[i].Type == 2) {
						action = " 回复了 "
					}

					if (messages[i].Sender == 0) {
						// 游客
						var item = '<div style="border-bottom: 1px solid #e2e2e2;min-width: 250px;padding-bottom: 3px;"><span style="color: #3e3e3e;">游客</span>' + action +
						'<a target="_blank" href="' + messages[i].Url + '?messageid=' + messages[i].Id + '#reply_id_' + messages[i].SourceId + '">' + messages[i].Title + '</a></div>';
						container.append(item);
					} else {
						var item = '<div style="border-bottom: 1px solid #e2e2e2;min-width: 250px;padding-bottom: 3px;"><a href="/member/' + messages[i].SenderName + '">' + messages[i].SenderName + '</a>' + action +
						'<a target="_blank" href="' + messages[i].Url + '?messageid=' + messages[i].Id + '#reply_id_' + messages[i].SourceId + '">' + messages[i].Title + '</a></div>';
						container.append(item);
					}
				}
				// Max tip
				if (totalMessageCnt > 8) {
					var item = '<div style="border-bottom: 1px solid #e2e2e2;min-width: 250px;padding-bottom: 3px;color:#a2a2a2;">评论太多了，请先阅读上面的吧...</div>';
					container.append(item);
				}
			}, 10)
		}
	})

	return '<div class="message_pop_container" id="id_message_pop_container"></div>';
}

$(document).ready(function(){
	pullMessageCount();
	adjustFooter();
	adjustBodyMinHeight();

	$("#navbar_message_popover").popover({
		html:true,
		title: "消息",
		content: formatMessageHTML
	});
})

$(window).resize(function(){
	adjustFooter();
	adjustBodyMinHeight();
})

$(document).resize(function(){
	adjustFooter();
	adjustBodyMinHeight();
});