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
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"
)

var parsedAssets map[string]map[AssetType]*asset

type AssetType byte

func init() {
	parsedAssets = map[string]map[AssetType]*asset{}
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
				Logger.Warning("%v \"%v\" can't find required file \"%v\"", this.assetType.String(), this.assetName, include_file)
			}
			this.include_files = append(this.include_files, file)
		}
	}
	this.build()
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
		if stat, err := os.Stat(file_path); !os.IsNotExist(err) {
			stat.ModTime()
			return file_path, nil
		}
	}
	return "", errors.New("Can't find file")
}
func (this *asset) assetExt() string {
	switch this.assetType {
	case ASSET_JAVASCRIPT : return ".js"
	case ASSET_STYLESHEET: return ".css"
	default: return ""
	}
}
func (this *asset) build() {
	if this.assetType == ASSET_STYLESHEET {
		if Config.MinifyCSS {

		}
		if Config.CombineCSS {

		}
		// minify css
	}
	if this.assetType == ASSET_JAVASCRIPT {
		if Config.MinifyJS {
			minifier := minify.New()
			for i, value := range this.include_files {
				filename := filepath.Base(value)
				filename = filename[0:len(filename) - 3]
				source_stat, err := os.Stat(value)
				if err != nil {
					Logger.Error("[ ASSET_PIPELINE ] Can't get info about source file '%s': %v", value, err)
					continue
				}
				if source_stat.IsDir() {
					Logger.Error("[ ASSET_PIPELINE ] Can't use directory in `require`: %s", value)
					continue
				}
				b_md5 := md5.Sum([]byte(source_stat.ModTime().String() + value))
				md5 := fmt.Sprintf("%x", b_md5)

				file, err := os.OpenFile(value, os.O_RDONLY, 0766)
				if err != nil {
					Logger.Error("[ ASSET_PIPELINE ] Can't open source file %v", value)
					continue
				}
				new_dir := filepath.Join(Config.TempDir, filepath.Dir(value), "/")
				err = os.MkdirAll(new_dir, 0766)
				if err != nil {
					Logger.Error("[ ASSET_PIPELINE ] Can't create temp dir: %v", err)
				}
				minified_path := filepath.Join(new_dir, filename + "-" + md5 + ".min.js")
				// If file already created-replace include files and ignore minifying step
				if _, err := os.Stat(minified_path); !os.IsNotExist(err) {
					this.include_files[i] = minified_path
					continue
				}
				minified_file, err := os.OpenFile(minified_path, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0766)
				if err != nil {
					Logger.Error("[ ASSET_PIPELINE ] Can't create new file: %v", err)
					continue
				}
				err = js.Minify(minifier, minified_file, file, map[string]string{})
				if err != nil {
					Logger.Error("[ ASSET_PIPELINE ] Minification error: %v", err)
					continue
				}
				this.include_files[i] = minified_path
			}
		}
		if Config.CombineJS {
			combined_md5 := md5.New()
			for _, src := range this.include_files {
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
			combined_path := filepath.Join(Config.TempDir, this.assetName + "-" + md5 + ".js")
			// If file already created-replace include files and ignore minifying step
			if _, err := os.Stat(combined_path); !os.IsNotExist(err) {
				this.include_files = []string{combined_path}
				return
			}
			combined_file, err := os.OpenFile(combined_path, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0766)
			if err != nil {
				Logger.Error("[ ASSET_PIPELINE ] Can't create new file: %v", err)
				return
			}
			for _, src := range this.include_files {
				file, err := os.OpenFile(src, os.O_RDONLY, 0766)
				if err != nil {
					Logger.Error("[ ASSET_PIPELINE ] Can't open source file %v", src)
					continue
				}
				file_reader := bufio.NewReader(file)
				file_reader.WriteTo(combined_file)
				combined_file.WriteString(";")
			}
			this.include_files = []string{combined_path}
		}
	}

}
func (this *asset) buildHTML() template.HTML {
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
	for _, value := range this.include_files {
		result += tag_fn("/" + value)
	}
	return result
}
func js_tag(location string) template.HTML {
	return template.HTML(fmt.Sprintf("script type=text/javascript src=\"%s\"></script>", location))
}
func css_tag(location string) template.HTML {
	return template.HTML(fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\"> ", location))
}
func parseAsset(assetName string, assetType AssetType) (*asset, error) {
	if v, ok := parsedAssets[assetName]; ok {
		if asset, ok := v[assetType]; ok {
			return asset, nil
		}
	}
	result := new(asset)
	result.assetType = assetType
	result.assetName = assetName
	err := result.parse()
	if err == nil {
		if _, ok := parsedAssets[assetName]; !ok {
			parsedAssets[assetName] = map[AssetType]*asset{}
		}
		parsedAssets[assetName][assetType] = result
	}
	return result, err
}
