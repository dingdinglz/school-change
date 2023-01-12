package server

import (
	"change/config"
	"change/logger"
	"github.com/gofiber/fiber/v2"
	WebLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"os"
)

// Start 启动服务器
func Start() {
	WebServer = fiber.New()
	WebServer.Use(WebLogger.New())
	WebServer.Static("/", "./web/public")
	WebServer.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("你好")
	})
	err := WebServer.Listen("0.0.0.0:" + config.ConfigData.Port)
	if err != nil {
		logger.ConsoleLogger.Errorln("服务器启动失败！", err.Error())
		logger.FileLogger.Errorln("服务器启动失败！", err.Error())
		os.Exit(1)
	}
}
