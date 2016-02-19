package js

import (
	"github.com/saturn4er/beego-assets"
	"fmt"
)

const JS_EXTENSION = ".js"

func init() {
	beego_assets.SetAssetFileExtension(JS_EXTENSION, beego_assets.ASSET_JAVASCRIPT)
	beego_assets.SetMinifyCallback(JS_EXTENSION, MinifyJavascript)
}

func MinifyJavascript(file *beego_assets.AssetFile) (result_path string, err error) {
	hash, err := file.GetHash()
	if err != nil {
		beego_assets.Logger.Error(err.Error())
		return
	}
	fmt.Println(hash)
	return "", nil
}