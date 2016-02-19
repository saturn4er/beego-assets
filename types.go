package beego_assets

import (
	"os"
)

const (
	ASSET_JAVASCRIPT AssetType = iota
	ASSET_STYLESHEET
)

type pre_afterBuildCallback func(asset *asset) error
type minifyFileCallback func(file *os.File) (result_file_path string, err error)

