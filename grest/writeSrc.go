package main

import (
	"fmt"
	"io/ioutil"
)

func writeMain(basePath, projectName string) error {

	src := fmt.Sprintf(`
package main

import (
	
	
	"flag"
	"log"
	"reflect"

	"github.com/dereking/grest" 
	"%s/controllers"   
)

func main() {
	conf := flag.String("conf", "app.conf", "the conf file in conf DIR for this server.")
	flag.Parse()

	log.Println("Starting server with config file :", *conf)

	s := grest.NewGrestServer(*conf)

	//controller register
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
    "log"
	"github.com/dereking/grest/mvc" 
	"github.com/dereking/grest/debug" 
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
	debug.Debug(arg)

	c.Session.Set("user", "ked")
	
	c.ViewBag["Title"] = arg.U
	c.ViewBag["cnt"] = 1024
	c.ViewBag["Msg"] = "你好." + arg.U
	c.ViewBag["Users"] = []string{"Jack", "Tomy", "James"}
 
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
	return ioutil.WriteFile(basePath+"controllers/HomeController.go", []byte(src), 0777)

}

func writeViewHome(basePath string) error {
	src := `<!DOCTYPE html>
<html lang="zh-CN">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!-- 上述3个meta标签*必须*放在最前面，任何其他内容都*必须*跟随其后！ -->
    <title>{{ .Title }}</title>

    <!-- Bootstrap -->
    <link href="/css/bootstrap.min.css" rel="stylesheet">

    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
      <script src="/js/html5shiv-3.7.2.min.js"></script>
      <script src="/js/respond-1.4.2.min.js"></script>
    <![endif]-->
  </head>
  <body>
	<div class="container">
		<div class="row">
			<div class="col-md-3"></div>
			<div class="col-md-6">
				<div class="panel panel-default">
					<div class="panel-heading">
						<h3 class="panel-title">REST Server</h3>
					</div>
					<div class="panel-body">
						<p>msg:{{ .Msg }} </p>
						<p>cnt {{ .cnt }}</p>
						
						{{range $k, $v := .Users}}
						    <div>{{$k}} => {{$v}} </div>  
						{{end}}
						
						<button class="btn btn-primary" onclick="alert('你好，世界！');">OK</button>
					</div> 
				</div> 
			</div>
			<div class="col-md-3"></div>
		</div>
	</div>

    <!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
    <script src="/js/jquery-1.12.2.min.js"></script>
    <!-- Include all compiled plugins (below), or include individual files as needed -->
    <script src="/js/bootstrap.min.js"></script>
  </body>
</html>`
	return ioutil.WriteFile(basePath+"views/Home/Index.html", []byte(src), 0777)

}

func writeModel(basePath string) error {
	src := "package models " +
		"type User struct {" +
		"	Name string `json:\"name\"`" +
		"	Age  int    `json:\"age\"`" +
		"}"
	return ioutil.WriteFile(basePath+"models/user.go", []byte(src), 0777)

}

func writeConf(basePath string) error {
	src := `#server working mode:  dev or prod
run = dev 

cache.expires=1h

cache.hosts=10.2.8.129:6379
cache.redis.password=root

cache.redis.maxidle=5
cache.redis.maxactive=0

#second
cache.redis.idletimeout=240
cache.redis.protocol=tcp

#ms
cache.redis.timeout.connect=10000 
cache.redis.timeout.read=5000
cache.redis.timeout.write=5000
 

db.mysql.hostWrite=10.2.8.129:3306
db.mysql.hostWrite.user=root
db.mysql.hostWrite.psw=root
db.mysql.hostWrite.dbName=db
db.mysql.hostWrite.maxOpenConns=200
db.mysql.hostWrite.maxIdleConns=100


db.mysql.hostRead=10.2.8.129:3306
db.mysql.hostRead.user=root
db.mysql.hostRead.psw=root
db.mysql.hostRead.dbName=db
db.mysql.hostRead.maxOpenConns=200
db.mysql.hostRead.maxIdleConns=100



# 允许访问的ip列表, 如果该项不存在, 则允许所有ip. 
allow.client.ip = 127.0.0.1;172.16.16.188;



[dev]
addr = 0.0.0.0:8000

[prod]
addr = 0.0.0.0:8000`
	return ioutil.WriteFile(basePath+"app.conf", []byte(src), 0777)

}

func writeReadme(basePath string) error {
	src := `
# 1.grest简介
> grest是为了快速开发rest api服务器而设计的一个web框架.

- 实现controller即可快速发布一个web服务.
- 主要目的是发布rest api
- 支持简单的mvc web页面. 没有做防xss等防护
- 支持 api 说明描述.
- 表单\query的字段不区分大小写.


# 安装grest
> 复制grest目录到$GOPATH/src目录下即可.

 

TODO: 开发创建项目的辅助工具. - [ ]

# 目录结构
- public 目录: 存放js\css\img等静态资源. 需发布
- conf 目录: 服务器配置文件存放路径. 需发布
- app 目录: 程序代码. 无需发布
- gorazor 目录: gohtml转换后go源码存放目录. 作为golang的包. 无需发布
- template 目录:存放gothml和静态html文件. 需发布
  - *.gohtml文件 不能在线修改. 更改后需要执行 gorazor 生成go源码, 然后编译整个工程
  - *.html 可以在线修改. 实时生效
  - layout 嵌套的模板页,必须是base.gohtml

# 如何编写模板 

 
# 引用

[^identity]:格力是珠海一家家电企业. http://www.gree.com
`
	return ioutil.WriteFile(basePath+"readme.md", []byte(src), 0777)

}
