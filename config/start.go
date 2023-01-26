package config

import (
	"change/tools"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// InitConfig 初始化配置
func InitConfig() {
	rootPath, _ := os.Getwd()
	if !tools.IsFileExist(filepath.Join(rootPath, "data", "change.config")) {
		configInitText, err := json.Marshal(&ConfigJsonModel{
			Port:   "80",
			SSL:    false,
			Debug:  false,
			School: "合肥市第七中学",
			Limit:  ConfigLimitModel{Change: 30},
		})
		if err != nil {
			fmt.Println("config创建出现错误！", err.Error())
			os.Exit(1)
		}
		_ = os.WriteFile(filepath.Join(rootPath, "data", "change.config"), configInitText, 0777)
	}
	configReadText, _ := os.ReadFile(filepath.Join(rootPath, "data", "change.config"))
	err := json.Unmarshal(configReadText, &ConfigData)
	if err != nil {
		fmt.Println("config解析出现错误！", err.Error())
		os.Exit(1)
	}
}
