[![Build Status](https://travis-ci.org/saturn4er/beego-assets.svg?branch=master)](https://travis-ci.org/saturn4er/beego-assets)
## Description
Rails-style assets for golang beego web framework

## Installation
To install library, create config in ./conf/asset-pipeline.conf and enter into the console:

	go get github.com/saturn4er/beego-assets

	
## Usage
Add following lines to your main file
 
	import (
		_ "github.com/saturn4er/beego-assets"
	)

Put basic config to ./conf/asset-pipeline.conf

Use functions javascript_include_tag and stylesheet_include_tag in your templates

	{{ javascript_include_tag "application" }}
	{{ stylesheet_include_tag "application" }}

## Less
Less using Node.js and node-less library, so, you should install it. 

After that you can use less by adding following line to your imports

		_ "github.com/saturn4er/beego-assets/less"

## Asset format
Asset extension sholud be .js or .css. Depends on include_tag function
Current version of library support only "require" method.

### CSS asset example
 
	/*= require css/file1
	/*= require css/file1
### JS asset example
 
	//= require js/file1
	//= require js/file2

## Config
The configuration file is ./conf/asset-pipeline.conf have basic INI format.

### Root params
- assets_dir - paths to assets files(You can specify many directories, separated by commas)
- public_dirs - paths, where library will search for files for assets.(You can specify many directories, separated by commas)
- temp_dir - path to store compiled asset files.

### Sections
You can define different parameters for different runmodes, which are defined in ./conf/app.conf. Name of section is the value of runmode. If there is no such section, all the parameters will be false
#### Parameters
- minify_js - Flag to minify javascript assets
- minify_css - Flag to minify stylesheet assets
- combine_js - Flag to combine multiple javascript files into one file
- combine_css - Flag to combine multiple stylesheet files into one file
- production_mode - It this flag is FALSE, assets will be recompiled each page request

### Config example:

	#paths where assets stored(Devided by comma)
	assets_dirs = assets/javascripts,assets/stylesheets
	#path to js/css
	public_dirs = static
	#where to put compiled assets files
	temp_dir = static/assets
	
	[production]
	minify_js = true
	minify_css = true
	
	combine_js = true
	combine_css = true
	
	production_mode = false

	[dev]
	minify_js = false
	minify_css = false
	
	combine_js = false
	combine_css = false
	
	production_mode = false
	
## API
- SetAssetFileExtension - add your own extension to files finder
	- extension string	- Extension name, with point.(Ex. ".js")
	- asset_type AssetType	- beego_assets.ASSET_STYLESHEET / beego_assets.ASSET_JAVASCRIPT
	
- SetPreBuildCallback  - define pre-build callback for assets. "cb" will be executed before building of asset
	- asset_type AssetType
	- cb pre_afterBuildCallback	- callback, which will be executed before asset compilation
	
- SetAfterBuildCallback  - define after-build callback for assets. "cb" will be executed after asset was built
	- asset_type AssetType
	- cb pre_afterBuildCallback	- callback, which will be executed after asset compilation
	
## Bugs
- When combine CSS files, neet to change relative paths.