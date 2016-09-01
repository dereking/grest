package mvc

//"github.com/dereking/grest/controller"

type IActionResult interface {
	ExecuteResult(c IController)
}

type ActionResult struct {
	HttpCode int
	Message  []byte
	Header   map[string]string
}

func NewActionResult(code int, msg string) *ActionResult {
	ret := &ActionResult{Header: make(map[string]string), HttpCode: code,
		Message: []byte(msg)}
	return ret
}

func (ar *ActionResult) ExecuteResult(c IController) {

	c.GetResponse().WriteHeader(ar.HttpCode)

	for k, v := range ar.Header {
		c.GetResponse().Header().Set(k, v)
	}

	c.GetResponse().Write(ar.Message)
}
