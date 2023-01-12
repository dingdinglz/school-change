package database

import (
	"change/logger"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"xorm.io/xorm"
)

// InitDatabase 初始化数据库
func InitDatabase() {
	rootPath, _ := os.Getwd()
	logger.ConsoleLogger.Infoln("正在启动数据库...")
	var err error
	DatabaseEngine, err = xorm.NewEngine("sqlite3", filepath.Join(rootPath, "data", "change.db"))
	if err != nil {
		logger.ConsoleLogger.Errorln("启动数据库失败！", err.Error())
		logger.FileLogger.Errorln("启动数据库失败！", err.Error())
		os.Exit(1)
	}
}
