//	多说评论
var duoshuoQuery = {short_name:"gocodecc"};
(function() {
	var ds = document.createElement('script');
	ds.type = 'text/javascript';ds.async = true;
	ds.src = (document.location.protocol == 'https:' ? 'https:' : 'http:') + '//static.duoshuo.com/embed.js';
	ds.charset = 'UTF-8';
	(document.getElementsByTagName('head')[0] 
		|| document.getElementsByTagName('body')[0]).appendChild(ds);
})();
	
//	调整尺寸

var timerHandle;
var dsDiv = $(".ds-thread");

function checkCommentLoaded() {
	var dsComments = dsDiv.find(".ds-comments");
	if (dsComments.length != 0) {
		adjustFooter();
		clearTimeout(timerHandle);
		timerHandle = null;
	}
}
timerHandle = setInterval(checkCommentLoaded, 10);

$(window).on('beforeunload',function(){
	if (null != timerHandle) {
		clearTimeout(timerHandle);
	}
});