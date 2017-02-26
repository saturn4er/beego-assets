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

const sassExtension = ".sass"
const sassExtensionLen = len(sassExtension)

var sassBuiltFilesDir = filepath.Join(beegoAssets.Config.TempDir, "sass")

func init() {
	_, err := exec.LookPath("sass")
	if err != nil {
		beegoAssets.Error("Please, install Node.js sass compiler: npm install node-sass -g")
		return
	}
	beegoAssets.SetAssetFileExtension(sassExtension, beegoAssets.AssetStylesheet)
	beegoAssets.SetPreLoadCallback(beegoAssets.AssetStylesheet, BuildSassAsset)
	err = os.MkdirAll(sassBuiltFilesDir, 0766)
	if err != nil {
		beegoAssets.Error(err.Error())
		return
	}
}

// BuildSassAsset - build sass file from beegoAssets.Asset
func BuildSassAsset(asset *beegoAssets.Asset) error {
	for i, src := range asset.IncludeFiles {
		ext := filepath.Ext(src)
		if ext == sassExtension {
			stat, err := os.Stat(src)
			if err != nil {
				beegoAssets.Error("Can't get stat of file %s. %v", src, err)
				continue
			}
			md5 := fmt.Sprintf("%x", md5.Sum([]byte(stat.ModTime().String()+src)))
			file := filepath.Base(src)
			fileName := file[:len(file)-sassExtensionLen]
			newFilePath := filepath.Join(sassBuiltFilesDir, fileName+"-"+md5+"_build.css")
			ex := exec.Command("sass", src, newFilePath)
			var out bytes.Buffer
			ex.Stderr = &out
			err = ex.Run()
			if err != nil {
				fmt.Println("Error building SASS file:")
				fmt.Println(out.String())
				continue
			}
			asset.IncludeFiles[i] = newFilePath
		}
	}
	return nil
}
