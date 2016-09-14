package mvc

import (
	"net/http"
	"time"
)

//"github.com/dereking/grest/controller"

type IActionResult interface {
	ExecuteResult(c IController)
}

type ActionResult struct {
	HttpCode int
	Message  []byte
	Header   map[string]string
	//Cookies http.CookieJar
}

func NewActionResult(code int, msg string) *ActionResult {
	ret := &ActionResult{
		Header:   make(map[string]string),
		HttpCode: code,
		Message:  []byte(msg),
		//Cookies: http.,
	}
	return ret
}

func (ar *ActionResult) ExecuteResult(c IController) {

	for k, v := range ar.Header {
		c.GetResponse().Header().Set(k, v)
	}

	c.GetResponse().WriteHeader(ar.HttpCode)
	c.GetResponse().Write(ar.Message)
}

func (ar *ActionResult) WriteCookie(c IController) {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "sessionid", Value: "abcd", Expires: expiration}

	http.SetCookie(c.GetResponse(), &cookie)
}
