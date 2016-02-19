package beego_assets

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"fmt"
	"strings"
)

type assetPipelineConfig struct {
	Runmode             string
	// Paths to assets
	AssetsLocations     []string
	// Paths to js/css files
	PublicDirs          []string
	// Path to store compiled assets
	TempDir             string

	// Flags
	MinifyCSS           bool
	MinifyJS            bool
	CombineCSS          bool
	CombineJS           bool

	// Association of extension->AssetType
	extensions          map[string]AssetType

	// callbacks
	preBuildCallbacks   map[AssetType]func(*asset) (error)
	minifyCallbacks     map[string]minifyFileCallback
	afterBuildCallbacks map[AssetType]func(*asset) (error)
}

func (this *assetPipelineConfig) Parse(filename string) {
	config, err := config.NewConfig("ini", filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	Config.Runmode = beego.AppConfig.DefaultString("runmode", "dev")
	locations := config.DefaultString("assets_dirs", "")
	Config.AssetsLocations = strings.Split(locations, ",")

	public_dirs := config.DefaultString("public_dirs", "")
	Config.PublicDirs = strings.Split(public_dirs, ",")

	Config.TempDir = config.DefaultString("temp_dir", "static/assets")

	runmode_params, err := config.GetSection(Config.Runmode)
	if err != nil {
		Logger.Warn("Can't get section \"%v\" from config asset-pipeline.conf. Using default params", Config.Runmode)
	}
	getBoolFromMap(&runmode_params, "minify_css", &Config.MinifyCSS, false)
	getBoolFromMap(&runmode_params, "minify_js", &Config.MinifyJS, false)
	getBoolFromMap(&runmode_params, "combine_css", &Config.CombineCSS, false)
	getBoolFromMap(&runmode_params, "combine_js", &Config.CombineJS, false)

}
func getBoolFromMap(array *map[string]string, key string, variable *bool, default_value bool) {
	if val, ok := (*array)[key]; ok {
		_val := strings.ToLower(val)
		*variable = _val == "true" || _val == "1"
	}else {
		*variable = default_value
	}
}
func init() {
	Config.Parse("./conf/asset-pipeline.conf")
	Config.extensions = map[string]AssetType{}
	Config.preBuildCallbacks = map[AssetType]func(*asset) (error){}
	Config.minifyCallbacks = map[string]minifyFileCallback{}
	Config.afterBuildCallbacks = map[AssetType]func(*asset) (error){}
}

var Config assetPipelineConfig

