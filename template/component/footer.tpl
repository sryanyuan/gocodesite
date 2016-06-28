{{define "footer"}}
<footer id="id_footer" class="footer navbar-fixed-bottom">
	<div class="container">
		<p>
			<ul class="footer-links">
				<li><a href="/about"><i class="fa fa-question" aria-hidden="true"></i></a>
				{{if ne .config.GithubAddress ""}}
				<li class="muted">&middot;</li>
				<li><a href="{{.config.GithubAddress}}" target="_blank"><i class="fa fa-github" aria-hidden="true"></i></a></li>
				{{end}}
				{{if ne .config.WeiboAddress ""}}
				<li class="muted">&middot;</li>
				<li><a href="{{.config.WeiboAddress}}" target="_blank"><i class="fa fa-weibo" aria-hidden="true"></i></a></li>
				{{end}}
			</ul>
		</p>
		<p>
			Build with {{.goversion}} · Based on <a href="http://getbootstrap.com/" target="_blank">bootstrap</a> · {{getprocesstime .requesttime}}
		</p>
		<p><i class="fa fa-copyright"></i>2016-2016 gocode.cc</p>
	</div>
</footer>
{{end}}