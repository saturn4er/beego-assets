package less

import (
	"github.com/gtforge/beego-assets"
	"path/filepath"
	"fmt"
	"os/exec"
	"os"
	"crypto/md5"
	"bytes"
)

const SASS_EXTENSION = ".sass"
const SASS_EXTENSION_LEN = len(SASS_EXTENSION)

var sass_built_files_dir = filepath.Join(beego_assets.Config.TempDir, "sass")

func init() {
	_, err := exec.LookPath("sass")
	if err != nil {
		beego_assets.Error("Please, install Node.js sass compiler: npm install node-sass -g")
		return
	}
	beego_assets.SetAssetFileExtension(SASS_EXTENSION, beego_assets.ASSET_STYLESHEET)
	beego_assets.SetPreLoadCallback(beego_assets.ASSET_STYLESHEET, BuildSassAsset)
	err = os.MkdirAll(sass_built_files_dir, 0766)
	if err != nil {
		beego_assets.Error(err.Error())
		return
	}
}
func BuildSassAsset(asset *beego_assets.Asset) error {
	for i, src := range asset.Include_files {
		ext := filepath.Ext(src)
		if ext == SASS_EXTENSION {
			stat, err := os.Stat(src)
			if err != nil {
				beego_assets.Error("Can't get stat of file %s. %v", src, err)
				continue
			}
			md5 := md5.Sum([]byte(stat.ModTime().String() + src))
			md5_s := fmt.Sprintf("%x", md5)
			file := filepath.Base(src)
			file_name := file[:len(file) - SASS_EXTENSION_LEN]

			new_file_path := filepath.Join(sass_built_files_dir, file_name + "-" + md5_s + "_build.css")
			ex := exec.Command("sass", src, new_file_path)
			var out bytes.Buffer
			ex.Stderr = &out
			err = ex.Run()
			if err != nil {
				fmt.Println("Error building SASS file:")
				fmt.Println(out.String())
				continue
			}
			asset.Include_files[i] = new_file_path

		}
	}
	return nil
}
