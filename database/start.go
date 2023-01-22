package database

import (
	"change/logger"
	"change/tools"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

// InitDatabase 初始化数据库
func InitDatabase() {
	rootPath, _ := os.Getwd()
	logger.ConsoleLogger.Infoln("正在启动数据库...")
	var err error
	var notInitBefore bool = false
	if !tools.IsFileExist(filepath.Join(rootPath, "data", "change.db")) {
		notInitBefore = true
	}
	DatabaseEngine, err = xorm.NewEngine("sqlite3", filepath.Join(rootPath, "data", "change.db"))
	if err != nil {
		logger.ConsoleLogger.Errorln("启动数据库失败！", err.Error())
		logger.FileLogger.Errorln("启动数据库失败！", err.Error())
		os.Exit(1)
	}
	DatabaseEngine.SetMapper(names.GonicMapper{})
	err = DatabaseEngine.Sync2(new(UserModel), new(SettingModel), new(ApplyModel), new(SubjectModel), new(ChangeModel))
	if err != nil {
		logger.ConsoleLogger.Errorln("同步数据库模型失败！", err.Error())
		logger.FileLogger.Errorln("同步数据库模型失败！", err.Error())
		os.Exit(1)
	}
	if notInitBefore {
		cnt, _ := DatabaseEngine.Table(new(UserModel)).Count()
		if cnt == 0 {
			logger.ConsoleLogger.Infoln("数据库：正在初始化user表")
			err = UserCreateNew("admin", "adminadmin", 3, "127.0.0.1", 3, 34, "管理员")
			if err != nil {
				logger.ConsoleLogger.Errorln("初始化用户失败！", err.Error())
				logger.FileLogger.Errorln("初始化用户失败！", err.Error())
				os.Exit(1)
			}
			logger.ConsoleLogger.Infoln("初始化用户完成！管理员用户：admin 密码：adminadmin")
		}
	}
}
