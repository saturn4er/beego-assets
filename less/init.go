package less

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gtforge/beego-assets"
)

const lessExtension = ".less"
const lessExtensionLen = len(lessExtension)

var lessBuiltFilesDir = filepath.Join(beegoAssets.Config.TempDir, "less")

func init() {
	_, err := exec.LookPath("lessc")
	if err != nil {
		beegoAssets.Error("Please, install Node.js less compiler: npm install less -g")
		return
	}

	beegoAssets.SetAssetFileExtension(lessExtension, beegoAssets.AssetStylesheet)
	beegoAssets.SetPreLoadCallback(beegoAssets.AssetStylesheet, BuildLessAsset)
	err = os.MkdirAll(lessBuiltFilesDir, 0766)
	if err != nil {
		beegoAssets.Error(err.Error())
		return
	}
}

// BuildLessAsset - build less file from beegoAssets.Asset
func BuildLessAsset(asset *beegoAssets.Asset) error {
	for i, src := range asset.IncludeFiles {
		ext := filepath.Ext(src)
		if ext == lessExtension {
			stat, err := os.Stat(src)
			if err != nil {
				beegoAssets.Error("Can't get stat of file %s. %v", src, err)
				continue
			}
			md5 := fmt.Sprintf("%x", md5.Sum([]byte(stat.ModTime().String()+src)))
			file := filepath.Base(src)
			fileName := file[:len(file)-lessExtensionLen]
			newFilePath := filepath.Join(lessBuiltFilesDir, fileName+"-"+md5+"_build.css")
			ex := exec.Command("lessc", src, newFilePath)
			var out bytes.Buffer
			ex.Stderr = &out
			err = ex.Run()
			if err != nil {
				fmt.Println("Error building LESS file:")
				fmt.Println(out.String())
				continue
			}
			asset.IncludeFiles[i] = newFilePath
		}
	}
	return nil
}
