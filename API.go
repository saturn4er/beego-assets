package beego_assets

import (
	"errors"
	"fmt"
	"crypto/md5"
	"os"
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
func GetAssetFileHash(file *os.File) (string, error) {
	stat, err := file.Stat()

	if err != nil {
		return "", errors.New(fmt.Sprintf("Can't get info about source file '%s': %v", file.Name(), err))
	}
	if stat.IsDir() {
		return "", errors.New(fmt.Sprintf("Can't use directory in `require`: %s", file.Name()))
	}
	b_md5 := md5.Sum([]byte(stat.ModTime().String() + file.Name()))
	md5 := fmt.Sprintf("%x", b_md5)
	return md5, nil
}