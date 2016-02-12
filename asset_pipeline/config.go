package asset_pipeline

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"fmt"
	"strings"
)

type assetPipelineConfig struct {
	Runmode    string
	MinifyCSS  bool
	MinifyJS   bool
	CombineCSS bool
	CombineJS  bool
}

func (this *assetPipelineConfig) Parse(filename string) {
	config, err := config.NewConfig("ini", filename)
	if err != nil {
		fmt.Println(err)
	}
	Config.Runmode = beego.AppConfig.DefaultString("runmode", "dev")

	runmode_params, err := config.GetSection(Config.Runmode)
	if err != nil {
		logger.Warn("Can't get section \"%v\" from config asset-pipeline.conf. Using default params", Config.Runmode)
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
}

var Config assetPipelineConfig

