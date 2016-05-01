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
	"regexp"
	"github.com/u007/beego-cache"
	"fmt"
)

var minifier = minify.New()
var parsedAssets = map[string]map[AssetType]*Asset{}

type Asset struct {
	assetName     string
	assetType     AssetType
	needMinify    bool
	needCombine   bool
	Include_files []string
	result        []assetFile
	cache					beego_cache.Cache
	assetPath     string
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
			Debug("asset: %s", line)
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
	// Debug("finding in %v", Config.AssetsLocations)
	for _, value := range Config.AssetsLocations {
		// Debug("included %v || %v", this.Include_files, this.result)
		file_path := path.Join(value, this.assetName) + this.assetExt()
		// Debug("include %s || %v || %v", file_path, this.Include_files, this.result))
		this.assetPath = file_path
		if _, err := os.Stat(file_path); !os.IsNotExist(err) {
			if _, err := os.Stat(file_path); !os.IsNotExist(err) {
				return file_path, nil
			}
			return file_path, nil
		}
	}
	return "", errors.New(fmt.Sprintf("Can't find asset: %v", this))
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
	if cbcks, ok := Config.preLoadCallbacks[this.assetType]; ok {
		for _, value := range cbcks {
			err := value(this)
			if err != nil {
				Error(err.Error())
				return
			}
		}
	}
	this.result = this.readAllIncludeFiles()
	
	if this.assetType == ASSET_STYLESHEET {
		this.replaceRelLinks()
	}
	
	if !this.needMinify && !this.needCombine {
		return
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

func copyFileContents(src, dst string) (err error, cotent string) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
	    return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
    return
	}
	err = out.Sync()
	return
}

func (this *Asset) minify() {
	var minifyHandler func(string) (string, error)
	if this.assetType == ASSET_JAVASCRIPT {
		minifyHandler = MinifyJavascript
	}else {
		minifyHandler = MinifyStylesheet
	}
	for i, assetFile := range this.result {
		file_name := filepath.Base(assetFile.Path)
		file_ext := filepath.Ext(file_name)
		file_name = file_name[:len(file_name) - len(file_ext)]
		file_hash := GetAssetFileHash(&assetFile.Body)
		minified_path := filepath.Join(Config.TempDir, file_name + "-" + file_hash + file_ext)
		minified_body, err := minifyHandler(assetFile.Body)
		if err != nil {
			Error(err.Error())
		} else {
			this.result[i].Path = minified_path
			this.result[i].Body = minified_body
		}
	}
}

func (this *Asset) combine() {
	var files_devider string
	if this.assetType == ASSET_JAVASCRIPT {
		files_devider = ";"
	}
	result := make([]string, len(this.result))
	for i, assetFile := range this.result {
		result[i] = assetFile.Body
	}
	s_result := strings.Join(result, files_devider)
	combined_path := filepath.Join(Config.TempDir, this.assetName + "-" + GetAssetFileHash(&s_result) + this.assetExt())
	this.result = []assetFile{assetFile{Path:combined_path, Body:s_result}}
}

// Read all files from Include files and return map[path_to_file]body
func (this *Asset) readAllIncludeFiles() []assetFile {
	result := []assetFile{}
	
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
		
		//create md5 version for js
		if this.assetType == ASSET_JAVASCRIPT {
			if !beego_cache.GetCache().FileChanged(path){
				Debug("asset not changed: %s", path)
				_, _, path, err = beego_cache.GetCache().FileCacheStat(path)
			} else {
				// file changed or is new
				file_name := filepath.Base(path)
				file_ext := filepath.Ext(file_name)
				file_name = file_name[:len(file_name) - len(file_ext)]
				file_hash := GetAssetFileHash(&file_body)
				new_file_path := filepath.Join(Config.TempDir, file_name + "-" + file_hash + file_ext)
				Debug("Writing destination %s", new_file_path)
				f, err := os.Create(new_file_path)
				if err != nil {
					Error("Error creating JS", err.Error())
					continue
				}
				_, err = f.WriteString(file_body)
				if err != nil {
					Error("Error copying JS", err.Error())
					continue
				}
				fi, err := os.Stat(path)
				if err != nil {
				  Warning("Unable to stat asset %s", path)// Could not obtain stat, handle error
					continue
				}
				
				beego_cache.GetCache().CacheFile(path, fi, new_file_path)
				path = new_file_path
			}
			
		} //is js
		
		result = append(result, assetFile{Path:path, Body:file_body})
		file.Close()
	}
	return result
}

func (this *Asset) writeResultToFiles() {
	for _, assetFile := range this.result {
		path_dir := filepath.Dir(assetFile.Path)
		err := os.MkdirAll(path_dir, 0766)
		if err != nil {
			Error(path_dir)
			continue
		}
		file, err := os.OpenFile(assetFile.Path, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0766)
		if err != nil {
			Error("Can't write result file: %v", err)
		}
		_, err = file.WriteString(assetFile.Body)
		if err != nil {
			Error("Can't write body to result file: %v", err)
		}
		file.Close()
	}
}

func (this *Asset) replaceRelLinks() {
	regex, err := regexp.Compile("url\\( *['\"]? *(\\/?(?:(?:\\.{1,2}|[a-zA-Z0-9-_.]+)\\/)*[a-zA-Z0-9-_.]+\\.[a-zA-Z0-9-_.]+)(?:\\?.+?)? *['\"]? *\\)")
	if err != nil {
		Error(err.Error())
		return
	}
	for i, assetFile := range this.result {
		file_dir := filepath.Dir(assetFile.Path)
		urls := regex.FindAllStringSubmatch(assetFile.Body, -1)
		for _, mathces := range urls {
			url_dir := mathces[1]
			abs_path := filepath.Join("/", file_dir, url_dir)
			assetFile.Body = strings.Replace(assetFile.Body, url_dir, abs_path, -1)
		}
		this.result[i] = assetFile
	}
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
	for _, assetFile := range this.result {
		result += tag_fn("/" + assetFile.Path)
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
