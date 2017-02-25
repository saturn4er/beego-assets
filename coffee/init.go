package coffee

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gtforge/beego-assets"
)

const coffeeExtension = ".coffee"
const coffeeExtensionLen = len(coffeeExtension)

func init() {
	_, err := exec.LookPath("coffee")
	if err != nil {
		beegoAssets.Error("Please, install Node.js coffee-script compiler: npm install coffee-script -g")
		return
	}

	beegoAssets.SetAssetFileExtension(coffeeExtension, beegoAssets.AssetJavascript)
	beegoAssets.SetPreLoadCallback(beegoAssets.AssetJavascript, BuildCoffeeAsset)
}

//BuildCoffeeAsset - build coffee script file from beegoAssets.Asset
func BuildCoffeeAsset(asset *beegoAssets.Asset) error {
	for i, src := range asset.IncludeFiles {
		ext := filepath.Ext(src)
		if ext == coffeeExtension {
			stat, err := os.Stat(src)
			if err != nil {
				beegoAssets.Error("Can't get stat of file %s. %v", src, err)
				continue
			}
			md5 := fmt.Sprintf("%x", md5.Sum([]byte(stat.ModTime().String()+src)))
			file := filepath.Base(src)
			fileName := file[:len(file)-coffeeExtensionLen]
			newFileDir := filepath.Join(beegoAssets.Config.TempDir, "coffee")
			newFileName := fileName + "-" + md5 + "_build.js"
			ex := exec.Command("coffee", "-c", "-o", newFileDir, src)
			var out bytes.Buffer
			ex.Stderr = &out
			err = ex.Run()
			if err != nil {
				fmt.Println("Error building Coffee file:")
				fmt.Println(out.String())
				continue
			}
			newFilePath := filepath.Join(newFileDir, newFileName)
			err = os.Rename(filepath.Join(newFileDir, fileName+".js"), newFilePath)
			if err != nil {
				beegoAssets.Error(err.Error())
			}
			asset.IncludeFiles[i] = newFilePath

		}
	}
	return nil
}
