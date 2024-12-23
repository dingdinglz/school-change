package main

import (
	"change/tools"
	"os"
	"path/filepath"
)

// InitPath 初始化文件夹路径等
func InitPath() {
	rootPath, _ := os.Getwd()
	tools.CreatePathIfNotExist(filepath.Join(rootPath, "log"))
	tools.CreatePathIfNotExist(filepath.Join(rootPath, "data"))
	tools.CreatePathIfNotExist(filepath.Join(rootPath, "data", "avatar"))
	tools.CreatePathIfNotExist(filepath.Join(rootPath, "data", "change_pic"))
	tools.CreatePathIfNotExist(filepath.Join(rootPath, "data", "comment"))
}
