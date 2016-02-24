package less

import (
	"github.com/saturn4er/beego-assets"
	"path/filepath"
	"fmt"
	"os/exec"
)

const LESS_EXTENSION = ".less"
const LESS_EXTENSION_LEN = len(LESS_EXTENSION)

func init() {
	beego_assets.SetAssetFileExtension(LESS_EXTENSION, beego_assets.ASSET_STYLESHEET)
	beego_assets.SetPreBuildCallback(beego_assets.ASSET_STYLESHEET, BuildLessAsset)
}
func BuildLessAsset(asset *beego_assets.Asset) error {
	for i, src := range asset.Include_files {
		ext := filepath.Ext(src)
		if ext == LESS_EXTENSION {
			file := filepath.Base(src)
			file_name := file[:len(file) - LESS_EXTENSION_LEN]
			new_file_path := filepath.Join(beego_assets.Config.TempDir, "less", file_name + "_build.css")
			ex := exec.Command("lessc", src, new_file_path)
			ex.Start()
			err := ex.Wait()
			if err != nil {
				fmt.Println(err)
			}
			asset.Include_files[i] = new_file_path

		}
	}
	return nil
}