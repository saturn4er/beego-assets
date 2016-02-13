package asset_pipeline

import (
	"os"
	"bufio"
	"io"
	"errors"
	"path"
	"strings"
	"html/template"
	"fmt"
)

type AssetType byte

const (
	ASSET_JAVASCRIPT AssetType = iota
	ASSET_STYLESCHEET
)

func (this *AssetType) String() string {
	switch *this {
	case ASSET_JAVASCRIPT:
		return "Javascript asset"
	case ASSET_STYLESCHEET:
		return "Css asset"
	default:
		return "Unknown asset type"
	}
}

type asset struct {
	assetName     string
	assetType     AssetType
	include_files []string
}

func (this *asset) parse() error {
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
		// TODO: fix this hardcode :)
		if strings.HasPrefix(line, "//= require ") {
			include_file := line[12:]
			file, err := this.findIncludeFilePath(include_file)
			if err != nil {
				logger.Warning("%v \"%v\" can't find required file \"%v\"", this.assetType.String(), this.assetName, include_file)
			}
			this.include_files = append(this.include_files, file)
		}
	}
	return nil
}
func (this *asset) findAssetPath() (string, error) {
	for _, value := range Config.AssetsLocations {
		file_path := path.Join(value, this.assetName) + this.assetExt()
		if _, err := os.Stat(file_path); !os.IsNotExist(err) {
			return file_path, nil
		}
	}
	return "", errors.New("Can't find asset ")
}
func (this *asset) findIncludeFilePath(file string) (string, error) {
	for _, value := range Config.PublicDirs {
		file_path := path.Join(value, file) + this.assetExt()
		if _, err := os.Stat(file_path); !os.IsNotExist(err) {
			return file_path, nil
		}
	}
	return "", errors.New("Can't find file")
}
func (this *asset) assetExt() string {
	switch this.assetType {
	case ASSET_JAVASCRIPT : return ".js"
	case ASSET_STYLESCHEET: return ".css"
	default: return ""
	}
}
func (this *asset) buildHTML() template.HTML {
	var tag_fn func(string) template.HTML
	switch this.assetType {
	case ASSET_JAVASCRIPT:
		tag_fn = js_tag
	case ASSET_STYLESCHEET:
		tag_fn = css_tag
	default:
		fmt.Println("Unknown asset type")
		return ""
	}
	var result template.HTML
	for _, value := range this.include_files {
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
func ParseAsset(assetName string, assetType AssetType) (*asset, error) {
	result := new(asset)
	result.assetType = assetType
	result.assetName = assetName
	err := result.parse()
	return result, err
}
