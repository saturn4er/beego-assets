package beegoAssets

import (
	"os"
)

const (
	//AssetJavascript type
	AssetJavascript assetsType = iota
	//AssetStylesheet type
	AssetStylesheet
)

type assetsType byte
type preLoadCallback func(*Asset) error
type preAfterBuildCallback func(result []assetFile, asset *Asset) error
type minifyFileCallback func(file *os.File) (result_file_path string, err error)
type assetFile struct {
	Path string
	Body string
}

func (a *assetsType) String() string {
	switch *a {
	case AssetJavascript:
		return "Javascript Asset"
	case AssetStylesheet:
		return "Css Asset"
	default:
		return "Unknown Asset type"
	}
}
