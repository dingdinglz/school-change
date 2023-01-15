package server

import (
	"change/config"
	"change/logger"
	"github.com/gofiber/fiber/v2"
	WebLogger "github.com/gofiber/fiber/v2/middleware/logger"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"
	"os"
)

// Start 启动服务器
func Start() {
	viewEngine := html.New("./web/views", ".html")
	viewEngineAddFuncs(viewEngine)
	viewEngine.Reload(true)
	WebServer = fiber.New(fiber.Config{Views: viewEngine})
	SessionStore = session.New()
	WebServer.Use(WebLogger.New(), recover2.New())
	WebServer.Static("/", "./web/public")
	BindRoutes()
	err := WebServer.Listen("0.0.0.0:" + config.ConfigData.Port)
	if err != nil {
		logger.ConsoleLogger.Errorln("服务器启动失败！", err.Error())
		logger.FileLogger.Errorln("服务器启动失败！", err.Error())
		os.Exit(1)
	}
}
