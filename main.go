package BeegoAssetPipeline

import (
	"github.com/astaxie/beego"
	. "github.com/saturn4er/beego-assets/asset_pipeline"
)

func init() {
	beego.AddFuncMap("javascript_include_tag", JavascriptIncludeTag)
	beego.AddFuncMap("stylesheet_include_tag", StyleSheetIncludeTag)
}
