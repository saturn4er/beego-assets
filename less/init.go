package less

import (
	"github.com/saturn4er/beego-assets"
	"path/filepath"
	"fmt"
	"os/exec"
	"os"
	"crypto/md5"
)

const LESS_EXTENSION = ".less"
const LESS_EXTENSION_LEN = len(LESS_EXTENSION)
var less_built_files_dir = filepath.Join(beego_assets.Config.TempDir, "less")
func init() {
	_, err := exec.LookPath("lessc")
	if err != nil {
		beego_assets.Error("Please, install Node.js less compiler: npm install less -g")
		return
	}

	beego_assets.SetAssetFileExtension(LESS_EXTENSION, beego_assets.ASSET_STYLESHEET)
	beego_assets.SetPreLoadCallback(beego_assets.ASSET_STYLESHEET, BuildLessAsset)
	err = os.MkdirAll(less_built_files_dir, 0766)
	if err != nil {
		beego_assets.Error(err.Error())
		return
	}
}
func BuildLessAsset(asset *beego_assets.Asset) error {
	for i, src := range asset.Include_files {
		ext := filepath.Ext(src)
		if ext == LESS_EXTENSION {
			stat, err := os.Stat(src)
			if err != nil {
				beego_assets.Error("Can't get stat of file %s. %v", src, err)
				continue
			}
			md5 := md5.Sum([]byte(stat.ModTime().String() + src))
			md5_s := fmt.Sprintf("%x", md5)
			file := filepath.Base(src)
			file_name := file[:len(file) - LESS_EXTENSION_LEN]
			new_file_path := filepath.Join(less_built_files_dir, file_name + "-" + md5_s + "_build.css")
			ex := exec.Command("lessc", src, new_file_path)
			ex.Start()
			err = ex.Wait()
			if err != nil {
				fmt.Println(err)
			}
			asset.Include_files[i] = new_file_path

		}
	}
	return nil
}