package mvc

import (
	"net/http"
)

type RedirectResult struct {
	ActionResult
	Permanent bool
	Url       string
}

func NewRedirectResult(code int, url string) *RedirectResult {
	ret := &RedirectResult{Url: url}
	ret.Header = make(map[string]string)
	ret.HttpCode = code
	return ret
}

func (ar *RedirectResult) ExecuteResult(c IController) {

	http.Redirect(c.GetResponse(), c.GetRequest(), ar.Url, 302)
}
