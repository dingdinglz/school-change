package server

import (
	"change/config"
	"github.com/gofiber/fiber/v2"
)

func MakeViewMap() map[string]interface{} {
	return map[string]interface{}{
		"School": config.ConfigData.School,
	}
}

// BindRoutes 绑定路由
func BindRoutes() {
	WebServer.Get("/", indexRoute)
}

func indexRoute(ctx *fiber.Ctx) error {
	return ctx.Render("index", MakeViewMap(), "layout/main")
}
