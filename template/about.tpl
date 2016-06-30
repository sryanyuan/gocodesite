{{define "Title"}}关于{{end}}
{{define "content"}}
<div id="id-content" class="row">
	<div class="col-md-8 col-md-offset-1">
		<ul class="breadcrumb">
			<li><a href="/"><i class="fa fa-home"></i> 首页</a></li>
			<li class="active">关于</li>
		</ul>
		<div class="panel panel-default">
			<div class="panel-heading">关于本站</div>
			<div class="panel-body">
				<p>本站基本就是一个写后端的程序猿写的前端，<span id="id-about-textnoob">菜鸟</span>一只。</p>
				<p>后端部分参照了<a href="http://golangtc.com" target="_blank">Golangtc</a>和<a href="http://studygolang.com" target="_blank">studygolang</a>，十分感谢前辈们的开源分享。</p>
				<p>后端采用了golang来编写，主要使用了gorilla/mux和自带的template。前端使用了bootstrap，对于后端来说，各种css太头疼了，用了bootstrap后感到终于不用头疼了。</p>
				<p>以前粗略的看过一些前端知识，只是感觉好复杂，要记得东西好多，每次都半途而废。这次硬着头皮上了，写着写着感觉貌似有点儿感觉了，看来凡事还是不能因难而退。</p>
				<br/><br/>
				<p>站长写过游戏客户端、服务端、服务器，这次又来鼓捣前端了，真是越学越觉得自己什么都不会。</p>
			</div>
			<div class="panel-footer">
			</div>
		</div>
	</div>
</div>
{{end}}