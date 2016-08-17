{{define "footer"}}
<footer id="id_footer" class="footer navbar-fixed-bottom1">
	<div class="container">
		<p>
			
		</p>
		<p>
			Build with {{.goversion}} · Based on <a href="http://getbootstrap.com/" target="_blank">bootstrap</a> · {{getProcessTime .requesttime}}
		</p>
		<p>
			<ul id="id-footer-links" class="footer-links">
				<li><a href="/about"><i class="fa fa-question" aria-hidden="true"></i></a>
				{{if ne .config.GithubAddress ""}}
				<li class="muted">&middot;</li>
				<li><a href="{{.config.GithubAddress}}" target="_blank"><i class="fa fa-github" aria-hidden="true"></i></a></li>
				{{end}}
				{{if ne .config.WeiboAddress ""}}
				<li class="muted">&middot;</li>
				<li><a href="{{.config.WeiboAddress}}" target="_blank"><i class="fa fa-weibo" aria-hidden="true"></i></a></li>
				{{end}}
				<li><i class="fa fa-smile-o"></i></li>
				<li><i class="fa fa-copyright"></i>2016-2016 gocode.cc</li>
			</ul>
		</p>
	</div>
</footer>
{{end}}