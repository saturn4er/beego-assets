package beego_assets

import (
	"fmt"
	"os"
	"crypto/md5"
	"errors"
)

const (
	ASSET_JAVASCRIPT AssetType = iota
	ASSET_STYLESHEET
)


type AssetFile struct {
	os.File
}

func (this AssetFile) GetHash() (string, error) {
	stat, err := this.Stat()

	if err != nil {
		return "", errors.New(fmt.Sprintf("Can't get info about source file '%s': %v", this.Name(), err))
	}
	if stat.IsDir() {
		return "", errors.New(fmt.Sprintf("Can't use directory in `require`: %s", this.Name()))
	}
	b_md5 := md5.Sum([]byte(stat.ModTime().String() + this.Name()))
	md5 := fmt.Sprintf("%x", b_md5)
	return md5, nil
}

type pre_afterBuildCallback func(asset *asset) (result_file_path string, err error)
type minifyFileCallback func(file *AssetFile) (result_file_path string, err error)

