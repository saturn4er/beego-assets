package beego_assets

import (
	"os"
	"bufio"
	"io"
	"errors"
	"path"
	"strings"
	"html/template"
	"fmt"
	"path/filepath"
	"crypto/md5"
)

var parsedAssets map[string]map[AssetType]*Asset

type AssetType byte

func init() {
	parsedAssets = map[string]map[AssetType]*Asset{}
}
func (this *AssetType) String() string {
	switch *this {
	case ASSET_JAVASCRIPT:
		return "Javascript asset"
	case ASSET_STYLESHEET:
		return "Css asset"
	default:
		return "Unknown asset type"
	}
}

type Asset struct {
	assetName     string
	assetType     AssetType
	Include_files []string
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
				Logger.Warning("%v \"%v\" can't find required file \"%v\"", this.assetType.String(), this.assetName, include_file)
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
	if cbcks, ok := Config.preBuildCallbacks[this.assetType]; ok {
		for _, value := range cbcks {
			err := value(this)
			if err != nil {
				Logger.Error("[ ASSET_PIPELINE ] %v", err)
				return
			}
		}
	}
	if (this.assetType == ASSET_JAVASCRIPT && Config.MinifyJS) || (this.assetType == ASSET_STYLESHEET && Config.MinifyCSS) {
		this.minify()
	}
	if (this.assetType == ASSET_JAVASCRIPT && Config.CombineJS) || (this.assetType == ASSET_STYLESHEET && Config.CombineCSS) {
		this.combine()
	}
	if cbcks, ok := Config.afterBuildCallbacks[this.assetType]; ok {
		for _, value := range cbcks {
			err := value(this)
			if err != nil {
				Logger.Error("[ ASSET_PIPELINE ] %v", err)
				return
			}
		}
	}
}
func (this *Asset)  minify() {
	for i, asset_file := range this.Include_files {
		extension := filepath.Ext(asset_file)
		file, err := os.OpenFile(asset_file, os.O_RDONLY, 0766)
		if err != nil {
			Logger.Error("[ ASSET_PIPELINE ] Can't open source file %v", asset_file)
			continue
		}
		if callback, ok := Config.minifyCallbacks[extension]; ok {
			new_path, err := callback(file)
			if err != nil {
				Logger.Error("[ ASSET_PIPELINE ] %s", err.Error())
			} else {
				this.Include_files[i] = new_path
			}
		}
		file.Close()
	}
}
// Return hash of asset, bases on modification time of included files
func (this *Asset) getHash() string {
	combined_md5 := md5.New()
	for _, src := range this.Include_files {
		stat, err := os.Stat(src)
		if err != nil {
			Logger.Error("[ ASSET_PIPELINE ] Can't get info about source file '%s': %v", src, err)
			continue
		}
		if stat.IsDir() {
			Logger.Error("[ ASSET_PIPELINE ] Can't use directory in `require`: %s", src)
			continue
		}
		combined_md5.Write([]byte(stat.ModTime().String() + src))

	}

	md5 := fmt.Sprintf("%x", combined_md5.Sum([]byte{}))
	return md5
}
func (this *Asset) combine() {
	hash := this.getHash()
	combined_path := filepath.Join(Config.TempDir, this.assetName + "-" + hash + this.assetExt())
	// If file already created-replace include files and ignore minifying step
	if _, err := os.Stat(combined_path); !os.IsNotExist(err) {
		this.Include_files = []string{combined_path}
		return
	}
	combined_file, err := os.OpenFile(combined_path, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0766)
	if err != nil {
		Logger.Error("[ ASSET_PIPELINE ] Can't create new file: %v", err)
		return
	}
	for _, src := range this.Include_files {
		file, err := os.OpenFile(src, os.O_RDONLY, 0766)
		if err != nil {
			Logger.Error("[ ASSET_PIPELINE ] Can't open source file %v", src)
			continue
		}
		file_reader := bufio.NewReader(file)
		file_reader.WriteTo(combined_file)
		combined_file.WriteString("\n")
	}
	this.Include_files = []string{combined_path}
}
func (this *Asset) buildHTML() template.HTML {
	var tag_fn func(string) template.HTML
	switch this.assetType {
	case ASSET_JAVASCRIPT:
		tag_fn = js_tag
	case ASSET_STYLESHEET:
		tag_fn = css_tag
	default:
		fmt.Println("Unknown asset type")
		return ""
	}
	var result template.HTML
	for _, value := range this.Include_files {
		result += tag_fn("/" + value)
	}
	return result
}
func js_tag(location string) template.HTML {
	return template.HTML(fmt.Sprintf("<script type=text/javascript src=\"%s\"></script>", location))
}
func css_tag(location string) template.HTML {
	return template.HTML(fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\"> ", location))
}
func parseAsset(assetName string, assetType AssetType) (*Asset, error) {
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
	err := result.parse()
	if err == nil {
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
