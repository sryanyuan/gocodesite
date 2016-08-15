<!DOCTYPE html>
<html lang="zh-CN">
  <head>
    <meta charset="utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <!-- 上述3个meta标签*必须*放在最前面，任何其他内容都*必须*跟随其后！ -->
    <title>{{template "Title" .}} - GoCode</title>
    <!-- Bootstrap -->
    <link href="/static/css/bootstrap.min.css" rel="stylesheet" />
    <!-- Bootstrap theme -->
    <link href="/static/css/bootstrap-theme.min.css" rel="stylesheet" />
	<!-- Font awesome -->
	<link href="/static/css/font-awesome.min.css" rel="stylesheet" />
	<!-- Custom css -->
    <link href="/static/css/base.css" rel="stylesheet" />
	{{template "importcss"}}
    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
      <script src="//cdn.bootcss.com/html5shiv/3.7.2/html5shiv.min.js"></script>
      <script src="//cdn.bootcss.com/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->
  </head>
  <body>
  {{template "navbar" .}}
  {{template "content" .}}
  {{template "footer" .}}
  <!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
  <script src="/static/js/jquery.min.js"></script> 
  <!-- Include all compiled plugins (below), or include individual files as needed -->
  <script src="/static/js/bootstrap.min.js"></script>
  <!-- Custom js -->
  <script src="/static/js/base.js"></script>
  {{template "importjs" .}}
  </body>
</html>
