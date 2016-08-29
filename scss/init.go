package less

import (
  // "github.com/saturn4er/beego-assets"
  "github.com/u007/beego-assets"
  "path/filepath"
  "fmt"
  "os/exec"
  "os"
  "io"
  "crypto/md5"
  "github.com/u007/beego-cache"
  "bytes"
)

const SCSS_EXTENSION = ".scss"
const SCSS_EXTENSION_LEN = len(SCSS_EXTENSION)

var scss_built_files_dir = filepath.Join(beego_assets.Config.TempDir, "scss")

func init() {
  _, err := exec.LookPath("node-sass")
  if err != nil {
    beego_assets.Error("Please, install Node.js sass compiler: npm install node-sass -g")
    return
  }
  beego_assets.SetAssetFileExtension(SCSS_EXTENSION, beego_assets.ASSET_STYLESHEET)
  beego_assets.SetPreLoadCallback(beego_assets.ASSET_STYLESHEET, BuildScssAsset)
  err = os.MkdirAll(scss_built_files_dir, 0766)
  if err != nil {
    beego_assets.Error(err.Error())
    return
  }
}
func BuildScssAsset(asset *beego_assets.Asset) error {
  for i, src := range asset.Include_files {
    ext := filepath.Ext(src)
    
    stat, err := os.Stat(src)
    if err != nil {
      beego_assets.Error("Can't get stat of file %s. %v", src, err)
      continue
    }
    
    md5 := md5.Sum([]byte(stat.ModTime().String() + src))
    md5_s := fmt.Sprintf("%x", md5)
    file := filepath.Base(src)
    file_name := file[:len(file) - SCSS_EXTENSION_LEN]
    new_file_path := filepath.Join(scss_built_files_dir, file_name + "-" + md5_s + "_build.css")
    
    if ext == SCSS_EXTENSION {
      
      if _, err := os.Stat(new_file_path); os.IsNotExist(err) {
        // beego_assets.Debug("exec: %s %s %s", "node-sass", src, new_file_path)
        ex := exec.Command("node-sass", src, new_file_path)
        var out bytes.Buffer
        ex.Stderr = &out
        err = ex.Run()
        if err != nil {
          fmt.Println("Error building SCSS file:")
          fmt.Println(out.String())
          continue
        }
      } else {
        beego_assets.Debug("skipping: %s", new_file_path)
      }
      asset.Include_files[i] = new_file_path
    } else if (ext == ".css") {
      
      if _, err := os.Stat(new_file_path); os.IsNotExist(err) {
        // plain copy
        in, err := os.Open(src)
        if err != nil {
          beego_assets.Warning("Unable to open source %s", src)
          continue
        }
        defer in.Close()
        out, err := os.Create(new_file_path)
        if err != nil {
          beego_assets.Warning("Unable to create destination %s", new_file_path)
          continue
        }
        defer func() {
          cerr := out.Close()
          if err == nil {
            err = cerr
          }
        }()
        if _, err = io.Copy(out, in); err != nil {
          beego_assets.Warning("Unable to copy %s", err.Error())
          continue
        }
        err = out.Sync()
        if (err != nil) {
          beego_assets.Warning("Unable to copy %s", src)
          continue
        }
      } else {
        beego_assets.Debug("skipping: %s", new_file_path)
      }
      asset.Include_files[i] = new_file_path
    }
    
    beego_cache.GetCache().CacheFile(src, stat, new_file_path)
  }
  return nil
}
