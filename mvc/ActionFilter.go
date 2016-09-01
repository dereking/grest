package mvc

//"github.com/dereking/grest/actionresult"

//ActionFilter , .
type IActionFilter interface {
	OnExecuting(*ActionExecutingContext)
	OnExecuted(*ActionExecutedContext)
	OnResultExecuting(*ResultExecutingContext)
	OnResultExecuted(*ResultExecutedContext)
}

type ActionExecutingContext struct {
	ActionParameters interface{}
	ActionName       string
	Result           IActionResult
}
type ActionExecutedContext struct {
}
type ResultExecutingContext struct {
}
type ResultExecutedContext struct {
}

func NewActionExecutingContext(actName string, args interface{}) *ActionExecutingContext {
	ret := &ActionExecutingContext{}
	ret.ActionName = actName
	ret.ActionParameters = args
	return ret
}
