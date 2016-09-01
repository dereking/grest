package mvc

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"runtime"
	"strings"

	"github.com/dereking/grest/config"
	"github.com/dereking/grest/debug"
	"github.com/dereking/grest/templateManager"
	"github.com/dereking/grest/utils"
)

//Controller , the basic controller class
type Controller struct {
	ViewBag  map[string]interface{}
	Response http.ResponseWriter
	Request  *http.Request
	//Session http.se
}

func (c *Controller) GetViewBag() map[string]interface{} {
	return c.ViewBag
}
func (c *Controller) GetResponse() http.ResponseWriter {
	return c.Response
}
func (c *Controller) GetRequest() *http.Request {
	return c.Request
}

func (c *Controller) Initialize(w http.ResponseWriter, r *http.Request) {
	c.ViewBag = make(map[string]interface{})
	c.Response = w
	c.Request = r
}

// ViewThisAction return the view of caller.
func (c *Controller) ViewThis() IActionResult {

	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)

	//获取此函数调用者的controller name和action name
	controllerName, actionName := fetchControllerActionName(f.Name())

	debug.Debug("ViewThisAction", controllerName, actionName)

	return c.View(controllerName, actionName)
}

func (c *Controller) View(controllerName, actionName string) IActionResult {

	log.Println("Controller.View : ViewBag=", c.ViewBag)

	ret := NewActionResult(200, "")
	var err error
	ret.Message, err = templateManager.Render(controllerName, actionName, c.ViewBag)
	if err != nil {
		ret.HttpCode = 404
		ret.Message = []byte(err.Error())
		return ret
	}
	//log.Println(string(content))
	return ret
}

//Redirect 302
func (c *Controller) Redirect(url string) IActionResult {
	ret := NewRedirectResult(302, url)
	ret.Header["Location"] = url
	//http.Redirect(c.Response,c.Request,url,302)
	return ret
}

func (c *Controller) HttpNotFound(msg string) IActionResult {
	return NewActionResult(404, msg+" Page Not found")
}
func (c *Controller) HttpForbidden() IActionResult {
	return NewActionResult(http.StatusForbidden, "Forbidden")
}

func (c *Controller) HttpInternalError(msg string) IActionResult {

	return NewActionResult(500, msg)
}

func (c *Controller) HttpHtml(html string) IActionResult {

	return NewActionResult(200, html)
}

func (c *Controller) HttpText(text string) IActionResult {

	return NewActionResult(200, text)
}

//客户端ip校验. 允许调用则返回nil, 不允许则返回err
func (c *Controller) ClientIPCheck(req *http.Request) error {
	ip := utils.GetClientIP(req)

	allow, found := config.AppConfig.String("allow.client.ip")

	//配置文件内没有ip限制, 那么允许.
	if !found {
		return nil
	}

	als := strings.Split(allow, ";")
	for _, a := range als {
		if a != "" && a == ip {
			return nil
		}
	}
	log.Println("Client IP not allowed", ip, allow)
	return errors.New("Client IP not allowed")
}

//构造jsonResult
func (c *Controller) JsonResult(o interface{}) *JsonResult {
	data, _ := json.Marshal(o)

	ret := &JsonResult{}
	ret.HttpCode = 200
	ret.Message = data
	return ret
}

//构造JsonAPIErrResult
func (c *Controller) JsonAPIErrResult(err error) *JsonResult {

	var jsonErr struct {
		R   int    `json:"r"`
		Msg string `json:"msg"`
	}
	jsonErr.R = 400
	jsonErr.Msg = err.Error()
	data, _ := json.Marshal(jsonErr)

	//log.Println("JsonAPIErrResult", jsonErr, data)

	ret := &JsonResult{}
	ret.HttpCode = 200
	ret.Message = data
	return ret
}
