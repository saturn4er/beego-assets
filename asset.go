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
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"fmt"
	"strconv"
	"time"
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
	cache					cache.Cache
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
			beego.Debug("asset:", line)
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
	// beego.Debug(fmt.Sprintf("finding in %v", Config.AssetsLocations))
	for _, value := range Config.AssetsLocations {
		// beego.Debug(fmt.Sprintf("included %v || %v", this.Include_files, this.result))
		file_path := path.Join(value, this.assetName) + this.assetExt()
		// beego.Debug(fmt.Sprintf("include %s || %v || %v", file_path, this.Include_files, this.result))
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

func (this *Asset) cacheInit() error {
	if this.cache == nil {
		bm, err := cache.NewCache("redis", `{"conn":":6379"}`)  
		this.cache = bm
		if err != nil{
			beego.Error("cache init error: ", err.Error())
			return fmt.Errorf("cache init error: ", err.Error())
		}
	}
	return nil
}

func (this *Asset) cache_name(name string) string {
	res := strings.Replace(name, "/", ":", -1)
	return res
}

func (this *Asset) cache_get(name string) string {
	this.cacheInit()
	result := fmt.Sprintf("%s", this.cache.Get(this.cache_name(name)))
	// beego.Debug("cache_get", name, fmt.Sprintf("%q", result))
	return result
}

func (this *Asset) cache_set(name string, value string, timeout time.Duration) error {
	this.cacheInit()
	err := this.cache.Put(this.cache_name(name), value, timeout)
	return err
}

func (this *Asset) cache_set_max(name string, value string) error {
	return this.cache_set(name, value, 60 * 60 * 24 * 365 * time.Second)
}

func (this *Asset) cache_exists(name string) bool {
	this.cacheInit()
	return this.cache.IsExist(this.cache_name(name))
}

func (this *Asset) fileChanged(path string) bool {
	cache_size, cache_time, cache_new_file, err := this.fileCacheStat(path)
	if err != nil {
		beego.Debug("filecachestat error:", err.Error())
		return true
	}
	
	fi, err := os.Stat(path)
	if err != nil {
	  Warning("Unable to stat asset %s", path)// Could not obtain stat, handle error
		return true
	}
	if fi.Size() != cache_size || fi.ModTime() != cache_time  {
		return true
	}
	
	//ensure file exists
	if _, err := os.Stat(cache_new_file); os.IsNotExist(err) {
		Warning("Cache file missing", cache_new_file)
		return true
	}
	
	return false
}

func (this *Asset) fileCacheStat(path string) (file_size int64, file_modtime time.Time, file_dest string, err error) {
	name := fmt.Sprintf("file_%s", path)
	if this.cache_exists(name) {
		// beego.Debug("Cache exists", name)
		cache := this.cache_get(name)
		res   := strings.Split(cache, "|")
		// beego.Debug(fmt.Sprintf("result: %q", res))
		file_size, err := strconv.ParseInt(res[0], 10, 64)
		if (err != nil) {
			Warning("fileCacheStat: can't get size from cache for %s", path)
		}
		// beego.Debug("cache:", cache)
		// https://gobyexample.com/time-formatting-parsing
		file_modtime, err := time.Parse(time.RFC3339, res[1]) //TODO
		if (err != nil) {
			Warning("fileCacheStat: mod time from cache: %s: %s", res[1], err.Error())
		}
		file_dest  := res[2]
		return file_size, file_modtime, file_dest, nil
	} else {
		beego.Debug("Cache missing", name)
		return 0, time.Time{}, "", fmt.Errorf("File missing: %f", path)
	}
}

func (this *Asset) cacheFile(path string, stat os.FileInfo, new_file_path string) {
	this.cacheInit()
	name := fmt.Sprintf("file_%s", path)
	err := this.cache_set_max(name, fmt.Sprintf("%d|%s|%s", stat.Size(), stat.ModTime().Format(time.RFC3339), new_file_path))
	if err != nil {
		beego.Debug("error set cacheFile", err.Error())
	}
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
			if !this.fileChanged(path){
				beego.Debug("asset not changed", path)
				_, _, path, err = this.fileCacheStat(path)
			} else {
				// file changed or is new
				file_name := filepath.Base(path)
				file_ext := filepath.Ext(file_name)
				file_name = file_name[:len(file_name) - len(file_ext)]
				file_hash := GetAssetFileHash(&file_body)
				new_file_path := filepath.Join(Config.TempDir, file_name + "-" + file_hash + file_ext)
				beego.Debug("Writing destination", new_file_path)
				f, err := os.Create(new_file_path)
				if err != nil {
					beego.Error("Error creating JS", err.Error())
					continue
				}
				_, err = f.WriteString(file_body)
				if err != nil {
					beego.Error("Error copying JS", err.Error())
					continue
				}
				fi, err := os.Stat(path)
				if err != nil {
				  Warning("Unable to stat asset %s", path)// Could not obtain stat, handle error
					continue
				}
				
				this.cacheFile(path, fi, new_file_path)
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
