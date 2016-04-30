package less

import (
	// "github.com/saturn4er/beego-assets"
	"beego-assets"
	"github.com/astaxie/beego"
	"path/filepath"
	"fmt"
	"os/exec"
	"os"
	"crypto/md5"
	"bytes"
)

const SCSS_EXTENSION = ".scss"
const SCSS_EXTENSION_LEN = len(SCSS_EXTENSION)

var scss_built_files_dir = filepath.Join(beego_assets.Config.TempDir, "scss")

func init() {
	_, err := exec.LookPath("node-sass")
	if err != nil {
		beego_assets.Error("Please, install Node.js sass compiler: npm install node-sass -g")
		return
	}
	beego_assets.SetAssetFileExtension(SCSS_EXTENSION, beego_assets.ASSET_STYLESHEET)
	beego_assets.SetPreLoadCallback(beego_assets.ASSET_STYLESHEET, BuildScssAsset)
	err = os.MkdirAll(scss_built_files_dir, 0766)
	if err != nil {
		beego_assets.Error(err.Error())
		return
	}
}
func BuildScssAsset(asset *beego_assets.Asset) error {
	beego.Debug(fmt.Sprintf("scsssssssss building %v", asset))
	for i, src := range asset.Include_files {
		ext := filepath.Ext(src)
		beego.Debug(fmt.Sprintf("included: %i: %s, ext: %s", i, src, ext))
		if ext == SCSS_EXTENSION {
			stat, err := os.Stat(src)
			if err != nil {
				beego_assets.Error("Can't get stat of file %s. %v", src, err)
				continue
			}
			md5 := md5.Sum([]byte(stat.ModTime().String() + src))
			md5_s := fmt.Sprintf("%x", md5)
			file := filepath.Base(src)
			file_name := file[:len(file) - SCSS_EXTENSION_LEN]
			new_file_path := filepath.Join(scss_built_files_dir, file_name + "-" + md5_s + "_build.css")
			ex := exec.Command("node-sass", "--scss", src, new_file_path)
			var out bytes.Buffer
			ex.Stderr = &out
			err = ex.Run()
			if err != nil {
				fmt.Println("Error building SCSS file:")
				fmt.Println(out.String())
				continue
			}
			asset.Include_files[i] = new_file_path

		}
	}
	return nil
}
