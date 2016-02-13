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
	asset, err := ParseAsset(asset_name, ASSET_JAVASCRIPT)
	if err != nil {
		logger.Error(err.Error())
		return ""
	}
	return asset.buildHTML()
}
func StyleSheetIncludeTag(asset_name string) template.HTML {
	asset, err := ParseAsset(asset_name, ASSET_STYLESCHEET)
	if err != nil {
		logger.Error(err.Error())
		return ""
	}
	return asset.buildHTML()
}
