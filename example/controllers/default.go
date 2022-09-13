package controllers

import (
	"github.com/beego/beego"
)

// MainController - MainController
type MainController struct {
	beego.Controller
}

// Get - Index Page
func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}
