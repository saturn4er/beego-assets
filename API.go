package beego_assets

import (
	"fmt"
	"crypto/md5"
)

func SetAssetFileExtension(extension string, asset_type AssetType) {
	Config.extensions[asset_type] = append(Config.extensions[asset_type], extension)
}

func SetPreBuildCallback(asset_type AssetType, cb pre_afterBuildCallback) {
	if _, ok := Config.preBuildCallbacks[asset_type]; !ok {
		Config.preBuildCallbacks[asset_type] = []pre_afterBuildCallback{}
	}
	Config.preBuildCallbacks[asset_type] = append(Config.preBuildCallbacks[asset_type], cb)
}
func SetAfterBuildCallback(asset_type AssetType, cb pre_afterBuildCallback) {
	if _, ok := Config.afterBuildCallbacks[asset_type]; !ok {
		Config.afterBuildCallbacks[asset_type] = []pre_afterBuildCallback{}
	}
	Config.afterBuildCallbacks[asset_type] = append(Config.afterBuildCallbacks[asset_type], cb)
}
func SetMinifyCallback(extension string, cb minifyFileCallback) {
	Config.minifyCallbacks[extension] = cb
}
func GetAssetFileHash(body *string) (string) {
	b_md5 := md5.Sum([]byte(*body))
	md5 := fmt.Sprintf("%x", b_md5)
	return md5
}