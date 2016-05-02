package beego_assets

import (
	"github.com/astaxie/beego"
	"fmt"
)

const PREFIX = "[ ASSET_PIPELINE ] "

func Debug(format string, v... interface{}) {
	beego.Debug(fmt.Sprintf(PREFIX + format, v...))
}
func Warning(format string, v... interface{}) {
	beego.Warning(fmt.Sprintf(PREFIX + format, v...))
}
func Error(format string, v... interface{}) {
	beego.Error(fmt.Sprintf(PREFIX + format, v...))
}
