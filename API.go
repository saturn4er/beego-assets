package beego_assets

func SetAssetFileExtension(extension string, asset_type AssetType) {
	Config.extensions[extension] = asset_type
}

func SetPreBuildCallback(asset_type AssetType, cb pre_afterBuildCallback) {

}
func SetAfterBuildCallback(asset_type AssetType, cb pre_afterBuildCallback) {

}
func SetMinifyCallback(extension string, cb minifyFileCallback) {
	Config.minifyCallbacks[extension] = cb
}