package mvc

import (
	"net/http"
)

//"grest/actionresult"

type IController interface {
	GetViewBag() map[string]interface{}
	GetResponse() http.ResponseWriter
	GetRequest() *http.Request
}
