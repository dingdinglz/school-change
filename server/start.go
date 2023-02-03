package server

import (
	"change/config"
	"change/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	WebLogger "github.com/gofiber/fiber/v2/middleware/logger"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"
	"os"
	"time"
)

// Start 启动服务器
func Start() {
	viewEngine := html.New("./web/views", ".html")
	viewEngineAddFuncs(viewEngine)
	viewEngine.Reload(true)
	WebServer = fiber.New(fiber.Config{Views: viewEngine,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			i := MakeViewMap(ctx)
			i["Error"] = err.Error()
			return ctx.Render("error", i, "layout/main")
		}})
	SessionStore = session.New()
	WebServer.Use(recover2.New())
	WebServer.Use(etag.New())
	WebServer.Use(WebLogger.New())
	WebServer.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        config.ConfigData.Limit.Web * 60,
		Expiration: 60 * time.Second,
		LimitReached: func(ctx *fiber.Ctx) error {
			return ctx.JSON(MakeApiResMap(false, "访问次数限制！"))
		},
	}))
	WebServer.Static("/", "./web/public")
	BindChatWebsocket()
	BindRoutes()
	err := WebServer.Listen("0.0.0.0:" + config.ConfigData.Port)
	if err != nil {
		logger.ConsoleLogger.Errorln("服务器启动失败！", err.Error())
		logger.FileLogger.Errorln("服务器启动失败！", err.Error())
		os.Exit(1)
	}
}
