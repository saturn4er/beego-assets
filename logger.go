package beegoAssets

import "github.com/astaxie/beego/logs"

var logger = logs.NewLogger(10000)

const prefix = "[ ASSET_PIPELINE ] "

//Warning - print Warning in log
func Warning(format string, v ...interface{}) {
	logger.Warning(prefix+format, v...)
}

//Error - print Error in log
func Error(format string, v ...interface{}) {
	logger.Error(prefix+format, v...)
}
