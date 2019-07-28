package controllers

import (
	"github.com/dereking/grest/log"
	"go.uber.org/zap"

	"github.com/dereking/grest/mvc"
)

type HomeController struct {
	mvc.Controller
}

func (c *HomeController) OnExecuting(a *mvc.ActionExecutingContext) {
	log.Logger().Info("HomeController OnExecuting",
		zap.Any("ActionParameters", a.ActionParameters))

	switch a.ActionName {
	case "Login":
	default:
		//If you want to check the user's access priveleges,
		//you can do it here.
		//if a.Result != nil, then the current action will not been executed.
		//a.Result = c.Redirect("/Home/Login")
		//a.Result = c.HttpForbidden()
	}
}

//
//
func (c *HomeController) Index(arg struct {
	U   string
	Cnt int
	Id  int
}) mvc.IActionResult {
	log.Logger().Debug("args", zap.Any("args", arg))

	users := []string{"Jack", "Tomy", "James"}

	c.Session.Set("user", "ked")

	c.ViewData["Title"] = arg.U + "session:user=" + c.Session.GetString("user")
	c.ViewData["Msg"] = users[2]
	c.ViewData["Users"] = users

	return c.ViewThis()
}

func (c *HomeController) Login() mvc.IActionResult {
	return c.View("Home", "Login")
}

func (c *HomeController) Test(arg struct {
	Id int
}) mvc.IActionResult {

	var dat struct {
		Users []string
		Id int
	}
	dat.Users = []string{"Jack", "Tomy", "James"}
	dat.Id = arg.Id
	return c.JsonResult(dat)
}
