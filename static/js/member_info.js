$(document).ready(function(){
	return;
	//	get comment count
	var articleGroup = $("#member-post-list-group");
	var articleCount = articleGroup.attr("articleCount");
	if (articleCount == 0) {
		return;
	}
	
	//	get all sub children
	var articles = articleGroup.find("span.member-post-reply-count");
	if (articles.length == 0){
		return;
	}
	var getUrl = "http://api.duoshuo.com/threads/counts.jsonp?short_name=gocodecc&threads=";
	$.each(articles, function(i, item){
		getUrl += $(item).attr("articleId");
		if (i != articles.length - 1) {
			getUrl += ",";
		}
	})
	
	getUrl += "&callback=?";
	
	$.getJSON(getUrl, function(ret){
		if (0 == ret.code) {
			$.each(ret.response, function(i, rsp){
				$("#member-post-reply-count-"+rsp.thread_key).html(rsp.comments);
			})
		}
	})
})