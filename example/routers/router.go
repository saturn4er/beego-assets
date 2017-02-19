package routers

import (
	"github.com/gtforge/beego-assets/example/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
