
# 1.grest简介
> grest是为了快速开发web服务而设计的一个web框架. 框架参考asp.net mvc。

- 实现controller即可快速发布一个web服务. 
- mvc  
- []支持 api 说明描述.
- 表单\query的字段不区分大小写.
- action 参数为匿名struct类型，所有字段需大写开头，以实现变量绑定。
- 数据绑定：request 的query、表单数据会转换成action的参数。


# 安装grest
> go get github.com/dereking/grest
> go install github.com/dereking/grest/grest

# 新建一个项目
> grest new subdir/projectName



# 目录结构
- static 目录: 存放js\css\img等静态资源. 需发布
- view 目录:存放模板文件和和静态html文件. 需发布
  - *.html 文件 不能在线修改.  

# 如何编写模板 
* 使用go  template引擎
* 模块函数
 - html 函数：输出html文本，不对其进行html编码。

 