{{define "Title"}}注册新用户{{end}}
{{define "content"}}
<div id="id-content" class="row">
  <div class="col-md-6 col-md-offset-3">
    <div class="panel panel-default">
      <div class="panel-heading">注册新用户</div>
      <div class="panel-body">
		{{if ne .signup_result ""}}
		<div class="alert alert-warning" role="alert">
          <strong>错误!</strong> {{.signup_result}}
		</div>
		{{end}}
        <form class="simple_form " novalidate="novalidate" id="new_user" action="/signup" accept-charset="UTF-8" method="post">
		  <input name="utf8" type="hidden" value="&#x2713;" />
		  <input type="hidden" name="authenticity_token" value="+/wrGAYYuda+veh7jrtK1yKRW1Rt7eSjnoD7IXmxdsrIZuDrkZ3ixbHnLQhWvoPOymGwNBRXS1b2PWuZyhhfMQ==" />
        
        <div class="form-group">
          <input type="email" class="form-control input-lg" placeholder="用户名" hint="谨慎填写，以后不可修改，建议用 Twitter ID。" name="user[login]" id="user_login" />
        </div>
        <div class="form-group">
          <input class="form-control input-lg" placeholder="名字" type="text" name="user[name]" id="user_name" />
        </div>
        <div class="form-group">
          <input type="email" class="form-control input-lg" placeholder="Email" name="user[email]" id="user_email" />
        </div>
        <div class="form-group">
          <div class="checkbox">
            <label for="user_email_public" class="checkbox"><input name="user[email_public]" type="hidden" value="0" /><input type="checkbox" value="1" checked="checked" name="user[email_public]" id="user_email_public" /> 公开 Email</label>
          </div>
        </div>
        <div class="form-group">
        <input class="form-control input-lg" placeholder="密码" type="password" name="user[password]" id="user_password" />
        </div>
        <div class="form-group">
        <input class="form-control input-lg" placeholder="确认密码" type="password" name="user[password_confirmation]" id="user_password_confirmation" />
        </div>
        <div class="form-group">
          <div class="input-group">
            <input class="form-control input-lg" placeholder="验证码" name="_rucaptcha" type="text" autocorrect="off" autocapitalize="off" pattern="[0-9a-z]*" maxlength="4" autocomplete="off" />
            <span class="input-group-addon input-group-captcha"><a class="rucaptcha-image-box" href="#"><img class="rucaptcha-image" src="https://ruby-china.org/rucaptcha/" alt="Rucaptcha" /></a></span>
          </div>
        </div>

        <div class="form-group">
          <input type="submit" name="commit" value="提交注册信息" class="btn btn-primary" data-disable-with="正在提交" />
        </div>
		
		<p>已有账号？请 <a class="btn btn-default" href="/signin">登录</a><p>
		</form>
	  </div>
    </div>
  </div>
</div>
{{end}}