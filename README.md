# grest
a GO lang REST &amp; web framework.

# install
> go get github.com/dereking/grest

> go install github.com/dereking/grest/grest

# start a new project
usage:
``` bash
  grest SUBCMD ARGS
```
 e.g.create a new GREST project in $GOPATH:

``` bash
   grest new projectName
```

The project will be created at $GOPATH/src/ProjectName

# static files
The directory named "static" is the place which storages the static files.
The subdirectories "css", "js", "images", "fonts" storage the corresponding static files.;

# controller 
there are one Filter Function in controller.
* OnExecuting Function

# websocket

``` go
func (c *WsController) Chat(ws *websocket.Conn) {

	defer ws.Close()

	var err error
	var str string

	for { 
		str = "hello, I'm server."

		if err = websocket.Message.Send(ws, str); err != nil {
			break
		} else {
			time.Sleep(time.Second * 2)
		}
	}
}
```


# template Function
* html 
	输出html代码. 对字符串进行html关键词\标签转义.

