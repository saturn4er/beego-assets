package less

import (
	"github.com/saturn4er/beego-assets"
	"path/filepath"
	"fmt"
	"os/exec"
	"os"
	"crypto/md5"
	"bytes"
)

const COFFEE_EXTENSION = ".coffee"
const COFFEE_EXTENSION_LEN = len(COFFEE_EXTENSION)

func init() {
	_, err := exec.LookPath("coffee")
	if err != nil {
		beego_assets.Error("Please, install Node.js coffee-script compiler: npm install coffee-script -g")
		return
	}

	beego_assets.SetAssetFileExtension(COFFEE_EXTENSION, beego_assets.ASSET_JAVASCRIPT)
	beego_assets.SetPreLoadCallback(beego_assets.ASSET_JAVASCRIPT, BuildLessAsset)
}
func BuildLessAsset(asset *beego_assets.Asset) error {
	for i, src := range asset.Include_files {
		ext := filepath.Ext(src)
		if ext == COFFEE_EXTENSION {
			stat, err := os.Stat(src)
			if err != nil {
				beego_assets.Error("Can't get stat of file %s. %v", src, err)
				continue
			}
			md5 := md5.Sum([]byte(stat.ModTime().String() + src))
			md5_s := fmt.Sprintf("%x", md5)
			file := filepath.Base(src)
			file_name := file[:len(file) - COFFEE_EXTENSION_LEN]
			new_file_dir := filepath.Join(beego_assets.Config.TempDir, "coffee");
			new_file_name := file_name + "-" + md5_s + "_build.js"
			ex := exec.Command("coffee", "-c", "-o", new_file_dir, src)
			var out bytes.Buffer
			ex.Stderr = &out
			err = ex.Run()
			if err != nil {
				fmt.Println("Error building Coffee file:")
				fmt.Println(out.String())
				continue
			}
			new_file_path := filepath.Join(new_file_dir, new_file_name)
			os.Rename(filepath.Join(new_file_dir, file_name+".js"), new_file_path)
			asset.Include_files[i] = new_file_path

		}
	}
	return nil
}