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

type pre_afterBuildCallback func(file_path AssetPath) (result_file_path string, err error)
type minifyFileCallback func(file_path AssetPath) (result_file_path string, err error)

type AssetPath string

func (this *AssetPath) GetHash() (string, error) {
	source_stat, err := os.Stat(this.String())
	if err != nil {
		return "", errors.New(fmt.Sprintf("Can't get info about source file '%s': %v", this, err))
	}
	if source_stat.IsDir() {
		return "", errors.New(fmt.Sprintf("Can't use directory in `require`: %s", *this))
	}
	b_md5 := md5.Sum([]byte(source_stat.ModTime().String() + this.String()))
	md5 := fmt.Sprintf("%x", b_md5)
	return md5, nil
}

func (this AssetPath) String() string {
	return string(this)
}