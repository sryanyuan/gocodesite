{{define "navbar"}}
<nav id="navbar" class="navbar navbar-default navbar-fixed-top" role="navigation">
    <div class="container">
        <div class="navbar-header">
            <button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-ex1-collapse">
                <span class="sr-only">Toggle navigation</span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
            </button>
            <a href="/" class="navbar-brand">
				<img src="/static/img/logo.png" style="margin-top: -9px;">
			</a>
        </div>
        <div class="collapse navbar-collapse navbar-ex1-collapse">
			<!--left bar-->
            <ul class="nav navbar-nav">
                <li {{if eq .active "home"}}class="active"{{end}}>
					<a href="/">主页</a>
				</li>
            </ul>
            
            <!--right bar-->
            <ul id="id-navbar-right" class="nav navbar-nav navbar-right">
				<li class="dropdown">
				{{if eq .user.Uid 0}}
					<a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true" aria-expanded="false"><i class="fa fa-user"></i> 游客 </i><i class="fa fa-caret-down"></i></a>
					<ul class="dropdown-menu">
						<li>
							<a href="/account/signin"><i class="fa fa-sign-in"></i>&nbsp;&nbsp;登陆</a>
						</li>
						<li>
							<a href="/account/signup"><i class="fa fa-level-up"></i>&nbsp;&nbsp;注册</a>
						</li>
					</ul>
				{{else}}
					<a href="#" id="id-navbar-avatar" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true" aria-expanded="false">
					{{if eq .user.Sex 0}}
						{{if eq .user.Avatar ""}}
						<img alt="{{.user.UserName}}" class="avatar img-rounded" height="32" src="{{.imgPrefix}}/male.png" width="32" />
						{{end}}
					{{else}}
						{{if eq .user.Avatar ""}}
						<img alt="{{.user.UserName}}" class="avatar img-rounded" height="32" src="{{.imgPrefix}}/female.png" width="32" />
						{{end}}
					{{end}}
					<i class="fa fa-caret-down"></i></a>
						<ul id="id_loginmenu" class="dropdown-menu">
							<li class="dropdown-header">欢迎 {{.user.UserName}}</li>
							<li role="separator" class="divider"></li>
							<li>
								<a href="#"><i class="fa fa-cog"></i>&nbsp;&nbsp;用户中心</a>
							</li>
							<li>
								<a href="#"><i class="fa fa-sign-out"></i>&nbsp;&nbsp;登出</a>
							</li>
						</ul>
				{{end}}
				</li>
			</ul>
			
			<!--search bar-->
			<ul class="nav navbar-nav navbar-right">
				<li class="nav-search hidden-xs hidden-sm">
					<form class="navbar-form form-search active" action="/search" method="GET">
						<div class="form-group">
							<input class="form-control" name="q" type="text" value="" placeholder="搜索" />
						</div>
						<i class="fa btn-search fa-search"></i>
					</form>
				</li>
			</ul>
            
        </div>
    </div>
</nav>
{{end}}