package beego_assets

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"html/template"
)

var Logger = logs.NewLogger(10000)

func init() {
	Logger.SetLogger("console", "")
	beego.AddFuncMap("javascript_include_tag", javascriptIncludeTag)
	beego.AddFuncMap("stylesheet_include_tag", styleSheetIncludeTag)
}

func javascriptIncludeTag(asset_name string) template.HTML {
	asset, err := parseAsset(asset_name, ASSET_JAVASCRIPT)
	if err != nil {
		Logger.Error(err.Error())
		return ""
	}
	return asset.buildHTML()
}
func styleSheetIncludeTag(asset_name string) template.HTML {
	asset, err := parseAsset(asset_name, ASSET_STYLESHEET)
	if err != nil {
		Logger.Error(err.Error())
		return ""
	}
	return asset.buildHTML()
}