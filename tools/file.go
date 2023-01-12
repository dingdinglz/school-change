package tools

import (
	"os"
)

// IsFileExist 判断文件或文件夹是否存在
func IsFileExist(_path string) bool {
	_, res := os.Stat(_path)
	return res == nil
}

// CreatePathIfNotExist 如果文件夹不存在则创建文件夹
func CreatePathIfNotExist(_path string) {
	if !IsFileExist(_path) {
		_ = os.Mkdir(_path, 0777)
	}
}
