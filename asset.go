package beegoAssets

import (
	"bufio"
	"errors"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tdewolff/minify"
)

var minifier = minify.New()
var parsedAssets = map[string]map[assetsType]*Asset{}

// Asset - Asset type can be any asset(js,css,sass etc)
type Asset struct {
	assetName    string
	assetType    assetsType
	needMinify   bool
	needCombine  bool
	IncludeFiles []string
	result       []assetFile
}

// Find Asset and parse it
func (a *Asset) parse() error {
	assetPath, err := a.findAssetPath()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(assetPath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	assetReader := bufio.NewReader(file)
	for {
		_line, _, err := assetReader.ReadLine()
		if err == io.EOF {
			break
		}
		line := string(_line)
		var prefix string
		if a.assetType == AssetJavascript {
			prefix = "//= require "
		} else {
			prefix = "/*= require "
		}
		if strings.HasPrefix(line, prefix) {
			includeFile := line[len(prefix):]
			file, err := a.findIncludeFilePath(includeFile)
			if err != nil {
				Warning("%v \"%v\" can't find required file \"%v\"", a.assetType.String(), a.assetName, includeFile)
				continue
			}
			a.IncludeFiles = append(a.IncludeFiles, file)
		}
	}
	return nil
}

// Search for Asset file in Config.AssetsLocations locations list
func (a *Asset) findAssetPath() (string, error) {
	for _, value := range Config.AssetsLocations {
		filePath := path.Join(value, a.assetName) + a.assetExt()
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			if _, err := os.Stat(filePath); !os.IsNotExist(err) {
				return filePath, nil
			}
			return filePath, nil
		}
	}
	return "", errors.New("Can't find Asset ")
}

// Search for Asset included file in Config.PublicDirs locations list.
func (a *Asset) findIncludeFilePath(file string) (string, error) {
	var extensions []string
	if val, ok := Config.extensions[a.assetType]; ok {
		extensions = val
	}
	for _, value := range Config.PublicDirs {
		for _, ext := range extensions {
			filePath := path.Join(value, file) + ext
			if stat, err := os.Stat(filePath); !os.IsNotExist(err) {
				stat.ModTime()
				return filePath, nil
			}
		}

	}
	return "", errors.New("Can't find file")
}

func (a *Asset) assetExt() string {
	switch a.assetType {
	case AssetJavascript:
		return ".js"
	case AssetStylesheet:
		return ".css"
	default:
		return ""
	}
}

func (a *Asset) build() {
	if cbcks, ok := Config.preLoadCallbacks[a.assetType]; ok {
		for _, value := range cbcks {
			err := value(a)
			if err != nil {
				Error(err.Error())
				return
			}
		}
	}
	a.result = a.readAllIncludeFiles()
	if !a.needCombine && !a.needMinify {
		return
	}
	if a.assetType == AssetStylesheet {
		a.replaceRelLinks()
	}
	if cbcks, ok := Config.preBuildCallbacks[a.assetType]; ok {
		for _, value := range cbcks {
			err := value(a.result, a)
			if err != nil {
				Error(err.Error())
				return
			}
		}
	}
	if a.needMinify {
		a.minify()
	}
	if a.needCombine {
		a.combine()
	}
	if cbcks, ok := Config.afterBuildCallbacks[a.assetType]; ok {
		for _, value := range cbcks {
			err := value(a.result, a)
			if err != nil {
				Error(err.Error())
				return
			}
		}
	}
	a.writeResultToFiles()
}

func (a *Asset) minify() {
	var minifyHandler func(string) (string, error)
	if a.assetType == AssetJavascript {
		minifyHandler = MinifyJavascript
	} else {
		minifyHandler = MinifyStylesheet
	}
	for i, assetFile := range a.result {
		fileName := filepath.Base(assetFile.Path)
		fileExt := filepath.Ext(fileName)
		fileName = fileName[:len(fileName)-len(fileExt)]
		fileHash := GetAssetFileHash(&assetFile.Body)
		minifiedPath := filepath.Join(Config.TempDir, fileName+"-"+fileHash+fileExt)
		minifiedBody, err := minifyHandler(assetFile.Body)
		if err != nil {
			Error(err.Error())
		} else {
			a.result[i].Path = minifiedPath
			a.result[i].Body = minifiedBody
		}
	}
}

func (a *Asset) combine() {
	var filesDevider string
	if a.assetType == AssetJavascript {
		filesDevider = ";"
	}
	result := make([]string, len(a.result))
	for i, assetFile := range a.result {
		result[i] = assetFile.Body
	}
	sResult := strings.Join(result, filesDevider)
	combinedPath := filepath.Join(Config.TempDir, a.assetName+"-"+GetAssetFileHash(&sResult)+a.assetExt())
	a.result = []assetFile{{Path: combinedPath, Body: sResult}}
}

// Read all files from Include files and return map[path_to_file]body
func (a *Asset) readAllIncludeFiles() []assetFile {
	result := []assetFile{}

	for _, path := range a.IncludeFiles {
		file, err := os.OpenFile(path, os.O_RDONLY, 0766)
		if err != nil {
			Warning("Can't open file %s. %v", path, err)
		}
		fileBody := ""
		fileReader := bufio.NewReader(file)
		for {
			line, _, err := fileReader.ReadLine()
			if err == io.EOF {
				break
			}
			fileBody += string(line) + "\n"
		}
		result = append(result, assetFile{Path: path, Body: fileBody})
		err = file.Close()
		if err != nil {
			Warning("Can't close file %s. %v", path, err)
		}
	}
	return result
}

func (a *Asset) writeResultToFiles() {
	for _, assetFile := range a.result {
		pathDir := filepath.Dir(assetFile.Path)
		err := os.MkdirAll(pathDir, 0766)
		if err != nil {
			Error(pathDir)
			continue
		}
		file, err := os.OpenFile(assetFile.Path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0766)
		if err != nil {
			Error("Can't write result file: %v", err)
		}
		_, err = file.WriteString(assetFile.Body)
		if err != nil {
			Error("Can't write body to result file: %v", err)
		}
		err = file.Close()
		if err != nil {
			Warning("Can't close file %s. %v", file.Name(), err)
		}
	}
}

func (a *Asset) replaceRelLinks() {
	regex, err := regexp.Compile("url\\( *['\"]? *(\\/?(?:(?:\\.{1,2}|[a-zA-Z0-9-_.]+)\\/)*[a-zA-Z0-9-_.]+\\.[a-zA-Z0-9-_.]+)(?:\\?.+?)? *['\"]? *\\)")
	if err != nil {
		Error(err.Error())
		return
	}
	for i, assetFile := range a.result {
		fileDir := filepath.Dir(assetFile.Path)
		urls := regex.FindAllStringSubmatch(assetFile.Body, -1)
		for _, matches := range urls {
			absPath := filepath.Join("/", fileDir, matches[1])
			assetFile.Body = strings.Replace(assetFile.Body, matches[1], absPath, -1)
		}
		a.result[i] = assetFile
	}
}

// Return html string of result
func (a *Asset) buildHTML() template.HTML {
	var tagFn func(string) template.HTML
	switch a.assetType {
	case AssetJavascript:
		tagFn = jsTag
	case AssetStylesheet:
		tagFn = cssTag
	default:
		Error("Unknown Asset type")
		return ""
	}
	var result template.HTML
	for _, assetFile := range a.result {
		result += tagFn("/" + assetFile.Path)
	}
	return result
}

func getAsset(assetName string, assetType assetsType) (*Asset, error) {
	// Check if a Asset already built ( only if production mode enabled )
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
	if assetType == AssetJavascript {
		if Config.MinifyJS {
			result.needMinify = true
		}
		if Config.CombineJS {
			result.needCombine = true
		}
	} else {
		if Config.MinifyCSS {
			result.needMinify = true
		}
		if Config.CombineCSS {
			result.needCombine = true
		}
	}
	err := result.parse()
	if err == nil {
		// Add Asset to cache ( only if production mode enabled )
		if Config.ProductionMode {
			if _, ok := parsedAssets[assetName]; !ok {
				parsedAssets[assetName] = map[assetsType]*Asset{}
			}
			parsedAssets[assetName][assetType] = result
		}
	}
	result.build()
	return result, err
}
