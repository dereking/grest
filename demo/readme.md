
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
