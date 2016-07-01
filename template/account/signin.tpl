{{define "Title"}}注册新用户{{end}}
{{define "importcss"}}{{end}}
{{define "importjs"}}
<script src="/static/js/signin.js"></script>
{{end}}
{{define "content"}}
<div id="id-content" class="row">
  <div class="col-sm-4 col-sm-offset-4">
    <div class="panel panel-default">
      <div class="panel-heading">登录</div>
      <div class="panel-body">
		<div id="id-signin-hint" class="alert alert-danger hidden" role="alert">
          <strong>错误!</strong><span id="id-signin-hinttext">Nothing</span>
		</div>
        <form id="id-form-signin" class="simple_form " novalidate="novalidate" id="new_user" action="/account/signin" accept-charset="UTF-8" method="post"><input name="utf8" type="hidden" value="&#x2713;" /><input type="hidden" name="authenticity_token" value="hVFHgkLJmc9/lZRLf31GJ4pqNZ3rAaJL0AAQgc6/0zK2y4xx1UzC3HDPUTineI8+Ypre/ZK7Db64vYA5fRb6yQ==" />
        <div class="form-group">
          <input type="email" class="form-control input-lg" placeholder="用户名 / Email" name="user[login]" id="user_login" />
        </div>
        <div class="form-group">
          <input type="password" class="form-control input-lg" placeholder="密码" name="user[password]" id="user_password" />
        </div>
        <div id="id-signin-captchaInputGroup" class="form-group">
          <div class="input-group">
			<input type="text" id="captchaSolution" name="captchaSolution" placeholder="请输入右侧验证码" />
			<img id="id-signin-captchaimg" src="/captcha/{{.captchaid}}.png" alt="验证码" title="看不清，点击" />
			<input type="hidden" id="id-signin-captchaIdHolder" name="captchaid" value="{{.captchaid}}">
          </div>
        </div>

        <div class="from-group checkbox">
            <label for="user_remember_me">
              <input name="user[remember_me]" type="hidden" value="0" />
			  <input type="checkbox" value="1" name="user[remember_me]" id="user_remember_me" /> 记住登录状态
            </label>
        </div>
        <div class="form-group">
            <input id="id-signin-submit" type="submit" name="commit" value="登录" class="btn btn-primary" data-disable-with="正在登录" />
        </div>
</form>      </div>
      <div class="panel-footer">
        
  <a href="/account/signup">注册</a>

  <a href="/account/forgotpassword">忘记了密码?</a>



      </div>
    </div>
  </div>
  <!--div class="col-md-3">
    <div class="panel panel-default">
      <div class="panel-heading">用其他平台的帐号登录</div>
      <ul class="list-group">
        <li class="list-group-item"><a class="btn btn-default btn-lg btn-block" href="/account/auth/github"><i class='fa fa-github'></i> GitHub</a> </li>
      </ul>
    </div>
  </div-->
</div>
{{end}}