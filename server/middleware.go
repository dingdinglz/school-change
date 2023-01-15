package server

import (
	"change/database"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// middleAdminRoute 后台鉴权中间件
func middleAdminRoute(ctx *fiber.Ctx) error {
	s, _ := SessionStore.Get(ctx)
	UserUsername := s.Get("user_username")
	UserPassword := s.Get("user_password")
	UserLevel := s.Get("user_level")
	UserID := s.Get("user_id")
	if UserUsername == nil || UserPassword == nil || UserLevel == nil {
		s.Set("message_warning", "请先登录！")
		_ = s.Save()
		return ctx.Redirect("/")
	} else {
		UserLevel_i, _ := strconv.Atoi(fmt.Sprintf("%s", UserLevel))
		UserID_i, _ := strconv.Atoi(fmt.Sprintf("%s", UserID))
		// logger.ConsoleLogger.Debugln(fmt.Sprintf("%s", UserUsername), fmt.Sprintf("%s", UserPassword), UserLevel_i)
		if database.UserCheck(fmt.Sprintf("%s", UserUsername), fmt.Sprintf("%s", UserPassword), UserLevel_i, UserID_i) {
			if UserLevel_i < 2 {
				s.Set("message_warning", "权限等级不足，禁止进入后台！")
				_ = s.Save()
				return ctx.Redirect("/")
			}
			return ctx.Next()
		} else {
			s.Delete("user_username")
			s.Delete("user_password")
			s.Delete("user_level")
			s.Delete("user_id")
			s.Set("message_warning", "账号信息过期，请重新登录！")
			_ = s.Save()
		}
	}
	return ctx.Redirect("/")
}
