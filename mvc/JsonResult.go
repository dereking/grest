package mvc

import (
	"fmt"
)

type JsonResult struct {
	ActionResult
}

func NewJsonResult(code int, msg []byte) *JsonResult {
	ret := &JsonResult{}
	ret.Header = make(map[string]string)
	ret.HttpCode = code
	ret.Message = msg
	return ret
}

func (ar *JsonResult) ExecuteResult(c IController) {

	for k, v := range ar.Header {
		c.GetResponse().Header().Set(k, v)
	}

	c.GetResponse().Header().Set("Content-Type", "application/json; charset=utf-8")

	c.GetResponse().WriteHeader(ar.HttpCode)

	fmt.Println("===========", ar.HttpCode, c.GetResponse().Header())

	c.GetResponse().Write(ar.Message)
}
