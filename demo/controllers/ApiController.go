package controllers

import (
	"log"

	"github.com/dereking/grest/demo/models"
	"github.com/dereking/grest/mvc"
)

type ApiController struct {
	mvc.APIController
}

func (c *ApiController) OnExecuting(a *mvc.ActionExecutingContext) {
	log.Println("ApiController OnExecuting")

	if c.ClientIPCheck() != nil {
		//停止执行后继action
		a.Result = c.HttpForbidden()
	}

}

func (c *ApiController) Index(arg struct {
	U   string
	Cnt int
	Id  int
}) mvc.IActionResult {
	var dat struct {
		Users []string
	}
	dat.Users = []string{"Jack", "Tomy", "James"}
	return c.JsonResult(dat)
}

func (c *ApiController) User_v1() mvc.IActionResult {

	u := &models.User{"Lily", 12}

	return c.JsonResult(u)
}
