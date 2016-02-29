package main

import (
	_ "github.com/saturn4er/beego-assets"
	_ "github.com/saturn4er/beego-assets/less"
	_ "github.com/saturn4er/beego-assets/sass"
	_ "github.com/saturn4er/beego-assets/scss"
	_ "github.com/saturn4er/beego-assets/coffeescript"
	_ "github.com/saturn4er/beego-assets/example/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

