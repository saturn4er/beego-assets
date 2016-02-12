package asset_pipeline

import (
	"github.com/astaxie/beego/logs"
	"html/template"
)

var logger = logs.NewLogger(10000)

func init() {
	logger.SetLogger("console", "")
}

func JavascriptIncludeTag(asset_name string) template.HTML {

	return "123"
}
func StyleSheetIncludeTag(asset_name string) template.HTML {
	return "234"
}
