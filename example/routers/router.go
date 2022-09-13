package routers

import (
	"github.com/beego/beego"
	"github.com/gtforge/beego-assets/example/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
