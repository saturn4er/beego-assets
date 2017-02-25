package beegoAssets

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
)

type assetPipelineConfig struct {
	Runmode string
	// Paths to assets
	AssetsLocations []string
	// Paths to js/css files
	PublicDirs []string
	// Path to store compiled assets
	TempDir string

	// Flags
	MinifyCSS      bool
	MinifyJS       bool
	CombineCSS     bool
	CombineJS      bool
	ProductionMode bool

	// Association of assetsType->Array of extensions
	extensions map[assetsType][]string

	// callbacks
	preLoadCallbacks    map[assetsType][]preLoadCallback
	preBuildCallbacks   map[assetsType][]preAfterBuildCallback
	minifyCallbacks     map[string]minifyFileCallback
	afterBuildCallbacks map[assetsType][]preAfterBuildCallback
}

// Config - assets pipeline configuration
var Config assetPipelineConfig

func init() {
	Config.Parse("./conf/Asset-pipeline.conf")
	Config.extensions = map[assetsType][]string{}
	Config.preBuildCallbacks = map[assetsType][]preAfterBuildCallback{}
	Config.minifyCallbacks = map[string]minifyFileCallback{}
	Config.afterBuildCallbacks = map[assetsType][]preAfterBuildCallback{}
	Config.preLoadCallbacks = map[assetsType][]preLoadCallback{}
}

func (a *assetPipelineConfig) Parse(filename string) {
	config, err := config.NewConfig("ini", filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	Config.Runmode = beego.AppConfig.DefaultString("runmode", "dev")
	locations := config.DefaultString("assets_dirs", "")
	Config.AssetsLocations = strings.Split(locations, ",")

	publicDirs := config.DefaultString("publicDirs", "")
	Config.PublicDirs = strings.Split(publicDirs, ",")
	Config.TempDir = config.DefaultString("temp_dir", "static/assets")

	runmodeParams, err := config.GetSection(Config.Runmode)
	if err != nil {
		logger.Warn("Can't get section \"%v\" from config Asset-pipeline.conf. Using default params", Config.Runmode)
	}
	getBoolFromMap(&runmodeParams, "minify_css", &Config.MinifyCSS, false)
	getBoolFromMap(&runmodeParams, "minify_js", &Config.MinifyJS, false)
	getBoolFromMap(&runmodeParams, "combine_css", &Config.CombineCSS, false)
	getBoolFromMap(&runmodeParams, "combine_js", &Config.CombineJS, false)
	getBoolFromMap(&runmodeParams, "production_mode", &Config.ProductionMode, false)

}
func getBoolFromMap(array *map[string]string, key string, variable *bool, defaultValue bool) {
	if val, ok := (*array)[key]; ok {
		_val := strings.ToLower(val)
		*variable = _val == "true" || _val == "1"
	} else {
		*variable = defaultValue
	}
}
