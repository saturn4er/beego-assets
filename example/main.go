package main

import (
	"github.com/astaxie/beego"
	_ "github.com/gtforge/beego-assets"
	_ "github.com/gtforge/beego-assets/coffee"
	_ "github.com/gtforge/beego-assets/example/routers"
	_ "github.com/gtforge/beego-assets/less"
	_ "github.com/gtforge/beego-assets/sass"
	_ "github.com/gtforge/beego-assets/scss"
)

func main() {
	beego.Run()
}
