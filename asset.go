package beego_assets

import (
	"os"
	"bufio"
	"io"
	"errors"
	"path"
	"strings"
	"html/template"
	"path/filepath"
	"github.com/tdewolff/minify"
)

var minifier = minify.New()
var parsedAssets = map[string]map[AssetType]*Asset{}

type Asset struct {
	assetName     string
	assetType     AssetType
	needMinify    bool
	needCombine   bool
	Include_files []string
	result        map[string]string
}

// Find asset and parse it
func (this *Asset) parse() error {
	assetPath, err := this.findAssetPath()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(assetPath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	asset_reader := bufio.NewReader(file)
	for {
		_line, _, err := asset_reader.ReadLine()
		if err == io.EOF {
			break
		}
		line := string(_line)
		var prefix string
		if this.assetType == ASSET_JAVASCRIPT {
			prefix = "//= require "
		}else {
			prefix = "/*= require "
		}
		if strings.HasPrefix(line, prefix) {
			include_file := line[len(prefix):]
			file, err := this.findIncludeFilePath(include_file)
			if err != nil {
				Warning("%v \"%v\" can't find required file \"%v\"", this.assetType.String(), this.assetName, include_file)
				continue
			}
			this.Include_files = append(this.Include_files, file)
		}
	}
	return nil
}

// Search for asset file in Config.AssetsLocations locations list
func (this *Asset) findAssetPath() (string, error) {
	for _, value := range Config.AssetsLocations {
		file_path := path.Join(value, this.assetName) + this.assetExt()
		if _, err := os.Stat(file_path); !os.IsNotExist(err) {
			if _, err := os.Stat(file_path); !os.IsNotExist(err) {
				return file_path, nil
			}
			return file_path, nil
		}
	}
	return "", errors.New("Can't find asset ")
}

// Search for asset included file in Config.PublicDirs locations list.
func (this *Asset) findIncludeFilePath(file string) (string, error) {
	var extensions []string;
	if val, ok := Config.extensions[this.assetType]; ok {
		extensions = val
	}
	for _, value := range Config.PublicDirs {
		for _, ext := range extensions {
			file_path := path.Join(value, file) + ext
			if stat, err := os.Stat(file_path); !os.IsNotExist(err) {
				stat.ModTime()
				return file_path, nil
			}
		}

	}
	return "", errors.New("Can't find file")
}

func (this *Asset) assetExt() string {
	switch this.assetType {
	case ASSET_JAVASCRIPT : return ".js"
	case ASSET_STYLESHEET: return ".css"
	default: return ""
	}
}

func (this *Asset) build() {
	this.result = this.readAllIncludeFiles()
	if !this.needCombine && ! this.needMinify {
		return
	}
	if this.assetType == ASSET_STYLESHEET {
		this.replaceRelLinks()
	}
	if cbcks, ok := Config.preBuildCallbacks[this.assetType]; ok {
		for _, value := range cbcks {
			err := value(this.result, this)
			if err != nil {
				Error(err.Error())
				return
			}
		}
	}
	if this.needMinify {
		this.minify()
	}
	if this.needCombine {
		this.combine()
	}
	if cbcks, ok := Config.afterBuildCallbacks[this.assetType]; ok {
		for _, value := range cbcks {
			err := value(this.result, this)
			if err != nil {
				Error(err.Error())
				return
			}
		}
	}
	this.writeResultToFiles()
}
func (this *Asset) minify() {
	var minifyHandler func(string) (string, error)
	if this.assetType == ASSET_JAVASCRIPT {
		minifyHandler = MinifyJavascript
	}else {
		minifyHandler = MinifyStylesheet
	}
	for path, body := range this.result {
		delete(this.result, path)
		file_name := filepath.Base(path)
		file_ext := filepath.Ext(file_name)
		file_name = file_name[:len(file_name) - len(file_ext)]
		file_hash := GetAssetFileHash(&body)
		minified_path := filepath.Join(Config.TempDir, file_name + "-" + file_hash + file_ext)
		minified_body, err := minifyHandler(body)
		if err != nil {
			Error(err.Error())
		} else {
			this.result[minified_path] = minified_body
		}
	}
}

func (this *Asset) combine() {
	var files_devider string
	if this.assetType == ASSET_JAVASCRIPT {
		files_devider = ";"
	}
	result := ""
	for _, body := range this.result {
		result += body + files_devider
	}
	combined_path := filepath.Join(Config.TempDir, this.assetName + "-" + GetAssetFileHash(&result) + this.assetExt())
	this.result = map[string]string{combined_path:result}
}

// Read all files from Include files and return map[path_to_file]body
func (this *Asset) readAllIncludeFiles() (map[string]string) {
	result := map[string]string{}

	for _, path := range this.Include_files {
		file, err := os.OpenFile(path, os.O_RDONLY, 0766)
		if err != nil {
			Warning("Can't open file %s. %v", path, err)
		}
		file_body := ""
		file_reader := bufio.NewReader(file)
		for {
			line, _, err := file_reader.ReadLine()
			if err == io.EOF {
				break
			}
			file_body += string(line) + "\n"
		}
		result[path] = file_body
		file.Close()
	}
	return result
}

func (this *Asset) writeResultToFiles() {
	for path, body := range this.result {
		path_dir := filepath.Dir(path)
		err := os.MkdirAll(path_dir, 0766)
		if err != nil {
			Error(path_dir)
			continue
		}
		file, err := os.OpenFile(path, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0766)
		if err != nil {
			Error("Can't write result file: %v", err)
		}
		_, err = file.WriteString(body)
		if err != nil {
			Error("Can't write body to result file: %v", err)
		}
		file.Close()
	}
}

func (this *Asset) replaceRelLinks() {

}

// Return html string of result
func (this *Asset) buildHTML() template.HTML {
	var tag_fn func(string) template.HTML
	switch this.assetType {
	case ASSET_JAVASCRIPT:
		tag_fn = js_tag
	case ASSET_STYLESHEET:
		tag_fn = css_tag
	default:
		Error("Unknown asset type")
		return ""
	}
	var result template.HTML
	for path, _ := range this.result {
		result += tag_fn("/" + path)
	}
	return result
}

func getAsset(assetName string, assetType AssetType) (*Asset, error) {
	// Check if this asset already built ( only if production mode enabled )
	if Config.ProductionMode {
		if v, ok := parsedAssets[assetName]; ok {
			if asset, ok := v[assetType]; ok {
				return asset, nil
			}
		}
	}
	result := new(Asset)
	result.assetType = assetType
	result.assetName = assetName
	if assetType == ASSET_JAVASCRIPT {
		if Config.MinifyJS {
			result.needMinify = true
		}
		if Config.CombineJS {
			result.needCombine = true
		}
	}else {
		if Config.MinifyCSS {
			result.needMinify = true
		}
		if Config.CombineCSS {
			result.needCombine = true
		}
	}
	err := result.parse()
	if err == nil {
		// Add asset to cache ( only if production mode enabled )
		if Config.ProductionMode {
			if _, ok := parsedAssets[assetName]; !ok {
				parsedAssets[assetName] = map[AssetType]*Asset{}
			}
			parsedAssets[assetName][assetType] = result
		}
	}
	result.build()
	return result, err
}
