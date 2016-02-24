package beego_assets

import "github.com/astaxie/beego/logs"

var Logger = logs.NewLogger(10000)

const PREFIX = "[ ASSET_PIPELINE ] "

func Warning(format string, v... interface{}) {
	Logger.Warning(PREFIX + format, v...)
}
func Error(format string, v... interface{}) {
	Logger.Error(PREFIX + format, v...)
}