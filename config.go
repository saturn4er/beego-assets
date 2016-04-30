package beego_assets

import (
	"github.com/astaxie/beego"
	// "github.com/astaxie/beego/config"
	"github.com/u007/go_config"
	// "github.com/go-ini/ini"
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
	ProductionMode      bool

	// Association of AssetType->Array of extensions
	extensions          map[AssetType][]string

	// callbacks
	preLoadCallbacks    map[AssetType][]preLoadCallback
	preBuildCallbacks   map[AssetType][]pre_afterBuildCallback
	minifyCallbacks     map[string]minifyFileCallback
	afterBuildCallbacks map[AssetType][]pre_afterBuildCallback
}

func (this *assetPipelineConfig) Parse(filename string) {
	// config, err := config.NewConfig("ini", filename)
	config, err := go_config.NewConfigLoader("ini", filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	Config.Runmode = beego.AppConfig.DefaultString("runmode", "dev")
	locations := config.String(Config.Runmode, "assets_dirs", "")
	beego.Debug(fmt.Sprintf("asset locations:: env: %s || %v || direct minify_js: %s", Config.Runmode, locations, 
		config.String(Config.Runmode, "minify_js", "")))
	Config.AssetsLocations = strings.Split(locations, ",")

	public_dirs := config.String(Config.Runmode, "public_dirs", "")
	Config.PublicDirs = strings.Split(public_dirs, ",")
	Config.TempDir = config.String(Config.Runmode, "temp_dir", "")

	Config.MinifyCSS = config.Boolean(Config.Runmode, "minify_css", false)
	Config.MinifyJS = config.Boolean(Config.Runmode, "minify_js", false)
	Config.CombineCSS = config.Boolean(Config.Runmode, "combine_css", false)
	Config.CombineJS = config.Boolean(Config.Runmode, "combine_js", false)
	Config.ProductionMode = config.Boolean(Config.Runmode, "production_mode", false)
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
	Config.extensions = map[AssetType][]string{}
	Config.preBuildCallbacks = map[AssetType][]pre_afterBuildCallback{}
	Config.minifyCallbacks = map[string]minifyFileCallback{}
	Config.afterBuildCallbacks = map[AssetType][]pre_afterBuildCallback{}
	Config.preLoadCallbacks = map[AssetType][]preLoadCallback{}
}

var Config assetPipelineConfig
