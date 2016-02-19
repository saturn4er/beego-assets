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

func MinifyJavascript(file_path beego_assets.AssetPath) (result_path string, err error) {
	hash, err := file_path.GetHash()
	if err != nil {
		beego_assets.Logger.Error(err.Error())
		return
	}
	fmt.Println(hash)
	return "", nil
}