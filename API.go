package beego_assets

func SetAssetFileExtension(extension string, asset_type AssetType) {
	Config.extensions[asset_type] = append(Config.extensions[asset_type], extension)
}

func SetPreBuildCallback(asset_type AssetType, cb pre_afterBuildCallback) {

}
func SetAfterBuildCallback(asset_type AssetType, cb pre_afterBuildCallback) {

}
func SetMinifyCallback(extension string, cb minifyFileCallback) {
	Config.minifyCallbacks[extension] = cb
}