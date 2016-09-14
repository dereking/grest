package controllers

import (
	"log"

	"github.com/dereking/grest/debug"
	"github.com/dereking/grest/mvc"
)

type HomeController struct {
	mvc.Controller
}

func (c *HomeController) OnExecuting(a *mvc.ActionExecutingContext) {
	log.Println("HomeController OnExecuting", a.ActionParameters)

	switch a.ActionName {
	case "Login":
	default:
		//assign the a.Result will stop the next action executing.
		//a.Result = c.Redirect("/Home/Login")
	}
}

//
//
func (c *HomeController) Index(arg struct {
	U   string
	Cnt int
	Id  int
}) mvc.IActionResult {
	debug.Debug(arg)

	user := []string{"Jack", "Tomy", "James"}

	c.ViewBag["Title"] = arg.U + "session:user=" + c.Session.GetString("user")
	c.ViewBag["Msg"] = user[2]
	c.ViewBag["Users"] = user

	return c.ViewThis()
}

func (c *HomeController) Login() mvc.IActionResult {
	return c.View("Home", "Login")
}

func (c *HomeController) Test(arg struct {
	U string
}) mvc.IActionResult {

	c.Session.Set("user", "ked")
	c.ViewBag["U"] = arg.U
	return c.View("Home", "Test")
}