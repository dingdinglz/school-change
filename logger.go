package main

import (
	"change/config"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

var (
	ConsoleLogger *logrus.Logger
	FileLogger    *logrus.Logger
)

// InitLogger 初始化日志器
func InitLogger() {
	rootPath, _ := os.Getwd()
	ConsoleLogger = logrus.New()
	ConsoleLogger.SetFormatter(&logrus.TextFormatter{DisableColors: true})
	ConsoleLogger.SetOutput(os.Stdout)
	fileout, _ := os.OpenFile(filepath.Join(rootPath, "log", time.Now().Format("20060102")+".log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	FileLogger = logrus.New()
	FileLogger.SetFormatter(&logrus.TextFormatter{DisableColors: true})
	FileLogger.SetOutput(fileout)
	if config.ConfigData.Debug {
		FileLogger.SetLevel(logrus.DebugLevel)
		ConsoleLogger.SetLevel(logrus.DebugLevel)
	} else {
		FileLogger.SetLevel(logrus.InfoLevel)
		ConsoleLogger.SetLevel(logrus.InfoLevel)
	}
}
