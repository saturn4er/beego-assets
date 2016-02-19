package css

import (
	"github.com/saturn4er/beego-assets"
	"os"
	"path/filepath"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"fmt"
	"errors"
)

const CSS_EXTENSION = ".css"
const CSS_EXTENSION_LEN = len(CSS_EXTENSION)

var minifier *minify.M

func init() {
	minifier = minify.New()
	beego_assets.SetAssetFileExtension(CSS_EXTENSION, beego_assets.ASSET_STYLESHEET)
	beego_assets.SetMinifyCallback(CSS_EXTENSION, MinifyJavascript)
}

func MinifyJavascript(file *os.File) (string, error) {
	file_path := file.Name()
	filename := filepath.Base(file_path)
	filename = filename[0:len(filename) - CSS_EXTENSION_LEN]

	hash, err := beego_assets.GetAssetFileHash(file)
	if err != nil {
		return "", err
	}
	new_dir := filepath.Join(beego_assets.Config.TempDir, filepath.Dir(file_path), "/")
	err = os.MkdirAll(new_dir, 0766)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Can't create temp dir: %v", err))
	}
	minified_path := filepath.Join(new_dir, filename + "-" + hash + ".min"+ CSS_EXTENSION)

	// If file already created-replace include files and ignore minifying step
	if _, err := os.Stat(minified_path); !os.IsNotExist(err) {
		return minified_path, nil
	}
	minified_file, err := os.OpenFile(minified_path, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0766)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Can't create new file: %v", err))
	}
	err = css.Minify(minifier, minified_file, file, map[string]string{})
	if err != nil {
		return "", errors.New(fmt.Sprintf("Minification error: %v", err))
	}
	return minified_path, nil
}