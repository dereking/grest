package mvc

import (
	//"fmt"
	"github.com/dereking/grest/log"
	"go.uber.org/zap"
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

	log.Logger().Debug("===========", zap.Int("HttpCode", ar.HttpCode), zap.Any("head", c.GetResponse().Header()))

	c.GetResponse().Write(ar.Message)
}
