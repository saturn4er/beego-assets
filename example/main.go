package main

import (
	_ "github.com/saturn4er/beego-assets"
	_ "github.com/saturn4er/beego-assets/less"
	_ "github.com/saturn4er/beego-assets/example/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

