{{define "Title"}}注册新用户{{end}}
{{define "content"}}
<div id="id-content" class="row">
  <div class="row">
  <div class="col-md-4 col-md-offset-4">
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
        <div class="form-group">
          <div class="input-group">
            <input class="form-control input-lg" placeholder="验证码" name="_rucaptcha" type="text" autocorrect="off" autocapitalize="off" pattern="[0-9a-z]*" maxlength="4" autocomplete="off" />
            <span class="input-group-addon input-group-captcha"><a class="rucaptcha-image-box" href="#"><img class="rucaptcha-image" src="https://ruby-china.org/rucaptcha/" alt="Rucaptcha" /></a></span>
          </div>
        </div>

        <div class="from-group checkbox">
            <label for="user_remember_me">
              <input name="user[remember_me]" type="hidden" value="0" /><input type="checkbox" value="1" name="user[remember_me]" id="user_remember_me" /> 记住登录状态
            </label>
        </div>
        <div class="form-group">
            <input id="id-signin-submit" type="submit" name="commit" value="登录" class="btn btn-primary btn-lg btn-block" data-disable-with="正在登录" />
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
</div>
{{end}}