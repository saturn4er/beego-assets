package routers

import (
	"github.com/saturn4er/beego-assets/example/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
