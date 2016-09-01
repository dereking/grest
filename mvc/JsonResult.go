package mvc

type JsonResult struct {
	ActionResult
}

func (ar *JsonResult) ExecuteResult(c IController) {

	c.GetResponse().WriteHeader(ar.HttpCode)

	for k, v := range ar.Header {
		c.GetResponse().Header().Set(k, v)
	}

	c.GetResponse().Write(ar.Message)
}
