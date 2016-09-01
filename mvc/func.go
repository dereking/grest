package mvc

import (
	"fmt"
	"runtime"
	"strings"

	//"github.com/dereking/grest/actionresult"
)

func HttpInternalError(err interface{}) IActionResult {

	var msg string

	switch err.(type) {
	case string:
		msg = err.(string)
	case error:
		msg = err.(error).Error() + "\nfunc call stack:\n"

		lines := 0          //调用堆栈的行数
		for i := 2; ; i++ { //从第二级开始.
			pc, file, line, ok := runtime.Caller(i)
			if ok {
				f := runtime.FuncForPC(pc)

				//只显示用户代码里的调用堆栈,过滤掉系统的.
				if !strings.HasPrefix(f.Name(), "runtime") && !strings.HasPrefix(f.Name(), "reflect") {
					//补足每行的前缀空格.
					for k := 0; k < lines; k++ {
						msg = msg + "  "
					}
					lines++ //调用堆栈的行数
					msg = msg + fmt.Sprintf("%s  %s:%d\n", f.Name(), file, line)
				}
			} else {
				break
			}
		}
	default:

		msg = "内部错误"
	}

	return NewActionResult(500, msg)
}

//NewController create a controller instance.
func NewController() *Controller {
	return &Controller{ViewBag: make(map[string]interface{})}
}

//parseControllerName  Get the name of Caller controller and action,
// return controller's name without "Controller" suffix,   and action's name
//eg. github.com/dereking/grest/demo/controllers.(*HomeController).Index
//    will return "Home","Index"
func fetchControllerActionName(actionFuncFullName string) (string, string) {
	p1 := strings.Index(actionFuncFullName, "(*")
	if p1 > 0 {
		p2 := strings.Index(actionFuncFullName, "Controller)")
		if p2 > 0 {
			return actionFuncFullName[p1+2 : p2], actionFuncFullName[p2+12:]
		} else {
			return "Home", "Index"
		}
	} else {
		return "Home", "Index"
	}
}
