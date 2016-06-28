{{define "navbar"}}
<nav class="navbar navbar-default navbar-fixed-top" role="navigation" id="navbar">
    <div class="container">
      <!--div class="navbar-header">
		<button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target=".navbar-collapse">
			<span class="sr-only">Toggle navigation</span>
		</button> 
		<a class="navbar-brand" href="/">Project name</a>
	  </div-->
	  <div class="navbar-collapse collapse">
		<ul class="nav navbar-nav navbar-left">
		  <li>
			<a href="/" class="navbar-brand">
				<img src="/static/img/logo.png" style="margin-top: -9px;">
			</a>
		  </li>
          <li {{if eq .active "home"}}class="active"{{end}}>
            <a href="/">主页</a>
          </li>
		</ul>
		<ul class="nav navbar-nav navbar-right">
          <li class="dropdown">
			{{if eq .user.Uid 1}}
            <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true"
            aria-expanded="false"><i class="fa fa-user"></i> 请登录 <i class="fa fa-angle-down"></i></a>
			{{else}}
			<a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true"
            aria-expanded="false">{{.user.UserName}} <i class="fa fa-caret-down"></i></a>
            <ul id="id_loginmenu" class="dropdown-menu">
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
      </div>
    </div>
 </nav>
{{end}}