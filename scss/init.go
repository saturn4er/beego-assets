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

const scssExtension = ".scss"
const scssExtensionLen = len(scssExtension)

var scssBuiltFilesDir = filepath.Join(beegoAssets.Config.TempDir, "scss")

func init() {
	_, err := exec.LookPath("sass")
	if err != nil {
		beegoAssets.Error("Please, install Node.js sass compiler: npm install node-sass -g")
		return
	}
	beegoAssets.SetAssetFileExtension(scssExtension, beegoAssets.AssetStylesheet)
	beegoAssets.SetPreLoadCallback(beegoAssets.AssetStylesheet, BuildScssAsset)
	err = os.MkdirAll(scssBuiltFilesDir, 0766)
	if err != nil {
		beegoAssets.Error(err.Error())
		return
	}
}

// BuildScssAsset - build scss file from beegoAssets.Asset
func BuildScssAsset(asset *beegoAssets.Asset) error {
	for i, src := range asset.IncludeFiles {
		ext := filepath.Ext(src)
		if ext == scssExtension {
			stat, err := os.Stat(src)
			if err != nil {
				beegoAssets.Error("Can't get stat of file %s. %v", src, err)
				continue
			}
			md5 := fmt.Sprintf("%x", md5.Sum([]byte(stat.ModTime().String()+src)))
			file := filepath.Base(src)
			fileName := file[:len(file)-scssExtensionLen]
			newFilePath := filepath.Join(scssBuiltFilesDir, fileName+"-"+md5+"_build.css")
			ex := exec.Command("sass", "--scss", src, newFilePath)
			var out bytes.Buffer
			ex.Stderr = &out
			err = ex.Run()
			if err != nil {
				fmt.Println("Error building SCSS file:")
				fmt.Println(out.String())
				continue
			}
			asset.IncludeFiles[i] = newFilePath
		}
	}
	return nil
}
