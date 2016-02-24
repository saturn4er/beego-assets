package beego_assets

import (
	"os"
)

type AssetType byte

func (this *AssetType) String() string {
	switch *this {
	case ASSET_JAVASCRIPT:
		return "Javascript asset"
	case ASSET_STYLESHEET:
		return "Css asset"
	default:
		return "Unknown asset type"
	}
}


const (
	ASSET_JAVASCRIPT AssetType = iota
	ASSET_STYLESHEET
)

const CSS_EXTENSION = ".css"
const CSS_EXTENSION_LEN = len(CSS_EXTENSION)
const JS_EXTENSION = ".js"
const JS_EXTENSION_LEN = len(JS_EXTENSION)

type preLoadCallback func(*Asset) error
type pre_afterBuildCallback func(result map[string]string, asset *Asset) error
type minifyFileCallback func(file *os.File) (result_file_path string, err error)
