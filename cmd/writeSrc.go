package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func writeMain(basePath, projectName string) error {

	src := fmt.Sprintf(`
package main

import (
	"flag"
	"reflect" 

	"%s/controllers"   

	"github.com/dereking/grest/log"
	"go.uber.org/zap" 
	"github.com/dereking/grest"  
)

func main() {
	conf := flag.String("conf", "app.conf", "the conf file in conf DIR for this server.")
	flag.Parse()

	log.Logger().Info("Starting server with config file :", zap.Any("conf", *conf))

	s := grest.NewGrestServer(*conf)

	//controller register
	s.AddController("api", reflect.TypeOf(controllers.ApiController{}))
	s.AddController("Home", reflect.TypeOf(controllers.HomeController{}))

	//main loop
	s.Serve()
}
	`, projectName)
	return ioutil.WriteFile(basePath+"main.go", []byte(src), 0777)

}

func writeController(basePath string) error {
	src := `package controllers

import (
	"github.com/dereking/grest/log"
	"go.uber.org/zap"

	"github.com/dereking/grest/mvc"  
)

type HomeController struct {
	mvc.Controller
 
}

func (c *HomeController) OnExecuting(a *mvc.ActionExecutingContext) {
	log.Println("HomeController OnExecuting", a.ActionParameters)

	switch a.ActionName {
	case "Login":
	default:
		//If you want to check the user's access priveleges, 
		//you can do it here.
		//if a.Result != nil, then the current action will not been executed.
		//a.Result = c.Redirect("/Home/Login")
		//a.Result = c.HttpForbidden()
	}
}

func (c *HomeController) Index(arg struct {
	U   string
	Cnt int
	Id  int
}) mvc.IActionResult { 
	log.Logger().Debug("args", zap.Any("args", arg))

	c.Session.Set("user", "ked")
	users := []string{"Jack", "Tomy", "James"}
	
	c.ViewData["Title"] = arg.U
	c.ViewData["cnt"] = 1024
	c.ViewData["Msg"] = "你好." + arg.U
	c.ViewData["Users"] = users
 
	return c.ViewThis() 
}

func (c *HomeController) Test(arg struct {
	Id int
}) mvc.IActionResult {

	var dat struct {
		Users []string
		Id int
	}
	dat.Users = []string{"Jack", "Tomy", "James"}
	dat.Id = arg.Id
	return c.JsonResult(dat)
}`
	return ioutil.WriteFile(
		fmt.Sprintf("%s%s%c%s", basePath, "controllers", os.PathSeparator, "HomeController.go"),
		[]byte(src), 0777)

}

func writeViewHome(basePath string) error {
	src := ` 
{{$title := "hello"}}
<div class="container">
	<div class="row">
		<div class="col-md-3"></div>
		<div class="col-md-6">
			<div class="panel panel-default">
				<div class="panel-heading">
					<h3 class="panel-title">GREST Server</h3>
				</div>
				<div class="panel-body">
					<p>all:{{ . }} </p>
					<p>msg:{{ .Msg }} </p>
					<p>cnt {{ .cnt }}</p>
					<p>$title {{ $title }}</p>
					{{range $k, $v := .Users}}
						<div>{{$k}} => {{$v}} </div>
					{{end}}

					<button class="btn btn-primary" onclick="alert('你好，世界！');">OK</button>
				</div>
			</div>
		</div>
		<div class="col-md-3"></div>
	</div>
</div>`
	return ioutil.WriteFile(
		fmt.Sprintf("%s%s%c%s%c%s", basePath, "views", os.PathSeparator, "Home", os.PathSeparator, "Index.html"),
		[]byte(src), 0777)

}

func writeViewShared(basePath string) error {
	src := `<!DOCTYPE html>
<html lang="zh-CN">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="/css/bootstrap.min.css" >
    <!-- 上述3个meta标签*必须*放在最前面，任何其他内容都*必须*跟随其后！ -->
    <title>{{ .Title }}</title>

    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
      <script src="/js/html5shiv-3.7.2.min.js"></script>
      <script src="/js/respond-1.4.2.min.js"></script>
    <![endif]-->
  </head>
  <body> 
  	{{ template "/shared/header.html" }}
    {{ @RenderBody() }}
    {{ template "/shared/footer.html" }}

    <!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
    <script src="/js/jquery-3.2.1.slim.min.js"></script>
    <script src="/js/popper.min.js"></script>
    <!-- Include all compiled plugins (below), or include individual files as needed -->
    <script src="/js/bootstrap.min.js"></script>
  </body>
</html> `
	return ioutil.WriteFile(
		fmt.Sprintf("%s%s%c%s%c%s", basePath, "views", os.PathSeparator, "Shared", os.PathSeparator, "_Layout.html"),
		[]byte(src), 0777)

}

func writeModel(basePath string) error {
	src := "package models " +
		"type User struct {" +
		"	Name string `json:\"name\"`" +
		"	Age  int    `json:\"age\"`" +
		"}"
	return ioutil.WriteFile(
		fmt.Sprintf("%s%s%c%s", basePath, "models", os.PathSeparator, "user.go"),
		[]byte(src), 0777)

}

func writeGoMod(basePath, projName string) error {
	src := fmt.Sprintf(`module %s`, projName)
	return ioutil.WriteFile(
		fmt.Sprintf("%s%s", basePath, "go.mod"),
		[]byte(src), 0777)

}

func writeConf(basePath string) error {
	src := `#server working mode:  [dev|prod]
run = dev


TemplateDir = views
TemplateExt = .html

cache.expires=1h

cache.hosts=127.0.0.1:6379
cache.redis.password=

cache.redis.maxidle=5
cache.redis.maxactive=0

#second
cache.redis.idletimeout=240
cache.redis.protocol=tcp

#ms
cache.redis.timeout.connect=10000 
cache.redis.timeout.read=5000
cache.redis.timeout.write=5000


#mysql
db.mysql.hostWrite=127.0.0.1:3306
db.mysql.hostWrite.user=greeg
db.mysql.hostWrite.psw=
db.mysql.hostWrite.dbName=test
db.mysql.hostWrite.maxOpenConns=200
db.mysql.hostWrite.maxIdleConns=100


db.mysql.hostRead=127.0.0.1:3306
db.mysql.hostRead.user=greeg
db.mysql.hostRead.psw=
db.mysql.hostRead.dbName=test
db.mysql.hostRead.maxOpenConns=200
db.mysql.hostRead.maxIdleConns=100



# 允许访问的ip列表, 如果该项不存在, 则允许所有ip.  127.0.0.1;172.16.16.188;
allow.client.ip = 127.0.0.3

 

[dev]
addr = 0.0.0.0:8000
# Auto reload the modified template from disk?
AutoReloadTemplate = true

[prod]
addr = 0.0.0.0:8000
# Auto reload the modified template from disk?
AutoReloadTemplate = false`
	return ioutil.WriteFile(basePath+"app.conf", []byte(src), 0777)

}

func writeReadme(basePath string) error {
	src := `https://github.com/dereking/grest`
	return ioutil.WriteFile(basePath+"readme.md", []byte(src), 0777)

}
