package routers

import (
	"github.com/astaxie/beego"
	"github.com/gtforge/beego-assets/example/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
