package main

import (
	_ "github.com/gtforge/beego-assets"
	_ "github.com/gtforge/beego-assets/less"
	_ "github.com/gtforge/beego-assets/sass"
	_ "github.com/gtforge/beego-assets/scss"
	_ "github.com/gtforge/beego-assets/coffee"
	_ "github.com/gtforge/beego-assets/example/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

