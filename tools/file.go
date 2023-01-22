package tools

import (
	"io/fs"
	"os"
	"path/filepath"
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

// GetNumberOfFiles 取目录下文件数
func GetNumberOfFiles(_path string) int {
	all := 0
	_ = filepath.Walk(_path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		all++
		return nil
	})
	return all
}
