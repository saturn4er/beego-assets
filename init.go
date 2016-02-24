package beego_assets

import (
	"github.com/astaxie/beego"
	"html/template"
)

func init() {
	Logger.SetLogger("console", "")
	beego.AddFuncMap("asset_js", getAssetHelper(ASSET_JAVASCRIPT))
	beego.AddFuncMap("asset_css", getAssetHelper(ASSET_STYLESHEET))
	Config.extensions[ASSET_JAVASCRIPT] = []string{".js"}
	Config.extensions[ASSET_STYLESHEET] = []string{".css"}

}
func getAssetHelper(Type AssetType) func(string) template.HTML {
	return func(asset_name string) template.HTML {
		asset, err := getAsset(asset_name, Type)
		if err != nil {
			Logger.Error(err.Error())
			return ""
		}
		return asset.buildHTML()
	}
}