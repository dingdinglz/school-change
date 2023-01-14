package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var (
	WebServer    *fiber.App
	SessionStore *session.Store
)
