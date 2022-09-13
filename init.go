package beegoAssets

import (
	"html/template"

	"github.com/beego/beego"
)

func init() {
	err := logger.SetLogger("console", "")
	if err != nil {
		Error(err.Error())
	}
	err = beego.AddFuncMap("asset_js", getAssetHelper(AssetJavascript))
	if err != nil {
		Error(err.Error())
	}
	err = beego.AddFuncMap("asset_css", getAssetHelper(AssetStylesheet))
	if err != nil {
		Error(err.Error())
	}
	Config.extensions[AssetJavascript] = []string{".js"}
	Config.extensions[AssetStylesheet] = []string{".css"}

}
func getAssetHelper(Type assetsType) func(string) template.HTML {
	return func(asset_name string) template.HTML {
		asset, err := getAsset(asset_name, Type)
		if err != nil {
			logger.Error(err.Error())
			return ""
		}
		return asset.buildHTML()
	}
}
