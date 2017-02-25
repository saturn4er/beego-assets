package beegoAssets

import (
	"crypto/md5"
	"fmt"
)

// SetAssetFileExtension - add new Asset type
func SetAssetFileExtension(extension string, assetType assetsType) {
	Config.extensions[assetType] = append(Config.extensions[assetType], extension)
}

// SetPreLoadCallback - set callback on preload assets
func SetPreLoadCallback(assetType assetsType, cb preLoadCallback) {
	if _, ok := Config.preLoadCallbacks[assetType]; !ok {
		Config.preLoadCallbacks[assetType] = []preLoadCallback{}
	}
	Config.preLoadCallbacks[assetType] = append(Config.preLoadCallbacks[assetType], cb)
}

// SetPreBuildCallback - set callback on PreBuild assets
func SetPreBuildCallback(assetType assetsType, cb preAfterBuildCallback) {
	if _, ok := Config.preBuildCallbacks[assetType]; !ok {
		Config.preBuildCallbacks[assetType] = []preAfterBuildCallback{}
	}
	Config.preBuildCallbacks[assetType] = append(Config.preBuildCallbacks[assetType], cb)
}

// SetAfterBuildCallback - set callback on after build assets
func SetAfterBuildCallback(assetType assetsType, cb preAfterBuildCallback) {
	if _, ok := Config.afterBuildCallbacks[assetType]; !ok {
		Config.afterBuildCallbacks[assetType] = []preAfterBuildCallback{}
	}
	Config.afterBuildCallbacks[assetType] = append(Config.afterBuildCallbacks[assetType], cb)
}

// SetMinifyCallback -
func SetMinifyCallback(extension string, cb minifyFileCallback) {
	Config.minifyCallbacks[extension] = cb
}

// GetAssetFileHash - get hash(fingerprint) file by body
func GetAssetFileHash(body *string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(*body)))
}
