package mvc

import (
	"net/http"
)

//"grest/actionresult"

type IController interface {
	GetViewData() map[string]interface{}
	GetResponse() http.ResponseWriter
	GetRequest() *http.Request
}
