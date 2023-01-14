package server

import (
	"change/config"
	"change/database"
	"change/tools"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
	"strconv"
)

func MakeViewMap(c *fiber.Ctx) map[string]interface{} {
	ResMap := map[string]interface{}{
		"School": config.ConfigData.School,
	}
	s, _ := SessionStore.Get(c)
	UserUsername := s.Get("user_username")
	UserPassword := s.Get("user_password")
	UserLevel := s.Get("user_level")
	UserID := s.Get("user_id")
	if UserUsername == nil || UserPassword == nil || UserLevel == nil {
		ResMap["User_login"] = false
	} else {
		UserLevel_i, _ := strconv.Atoi(fmt.Sprintf("%s", UserLevel))
		UserID_i, _ := strconv.Atoi(fmt.Sprintf("%s", UserID))
		// logger.ConsoleLogger.Debugln(fmt.Sprintf("%s", UserUsername), fmt.Sprintf("%s", UserPassword), UserLevel_i)
		if database.UserCheck(fmt.Sprintf("%s", UserUsername), fmt.Sprintf("%s", UserPassword), UserLevel_i, UserID_i) {
			ResMap["User_login"] = true
			ResMap["User_username"] = fmt.Sprintf("%s", UserUsername)
			ResMap["User_level"] = UserLevel_i
			ResMap["User_id"] = UserID_i
		} else {
			ResMap["User_login"] = false
			ResMap["WarningMessage"] = "登录信息有误！请重新登录！"
			s.Delete("user_username")
			s.Delete("user_password")
			s.Delete("user_level")
			s.Delete("user_id")
			_ = s.Save()
		}
	}
	return ResMap
}

// BindRoutes 绑定路由
func BindRoutes() {
	WebServer.Get("/", indexRoute)
	WebServer.Get("/avatar/:username", avatarRoute)
	apiRoute := WebServer.Group("/api")
	apiRoute.Post("/login", apiLoginRoute)
	apiRoute.Get("/logout", apiLogoutRoute)
}

// indexRoute 主页路由
func indexRoute(ctx *fiber.Ctx) error {
	return ctx.Render("index", MakeViewMap(ctx), "layout/main")
}

// MakeApiResMap 生成返回的json
func MakeApiResMap(ok bool, message string) fiber.Map {
	if ok {
		return fiber.Map{"status": "ok", "message": message}
	}
	return fiber.Map{"status": "err", "message": message}
}

// apiLoginRoute 登录api
func apiLoginRoute(ctx *fiber.Ctx) error {
	username := ctx.FormValue("username", "")
	password := ctx.FormValue("password", "")
	if username == "" || password == "" {
		return ctx.JSON(MakeApiResMap(false, "账号或密码为空！"))
	}
	if !database.UserHaveUserByUsername(username) {
		return ctx.JSON(MakeApiResMap(false, "用户不存在！"))
	}
	u := database.UserGetUserByUsername(username)
	if u.Password != tools.MD5(password) {
		return ctx.JSON(MakeApiResMap(false, "密码错误！"))
	}
	s, _ := SessionStore.Get(ctx)
	s.Set("user_username", u.Username)
	s.Set("user_password", u.Password)
	s.Set("user_level", strconv.Itoa(u.Level))
	s.Set("user_id", strconv.Itoa(u.ID))
	_ = s.Save()
	return ctx.JSON(MakeApiResMap(true, "登录成功！"))
}

// avatarRoute 头像路由
func avatarRoute(ctx *fiber.Ctx) error {
	username := ctx.Params("username", "")
	rootPath, _ := os.Getwd()
	if !tools.IsFileExist(filepath.Join(rootPath, "data", "avatar", username+".png")) {
		return ctx.SendFile(filepath.Join(rootPath, "web", "picture", "people.png"))
	}
	return ctx.SendFile(filepath.Join(rootPath, "data", "avatar", username+".png"))
}

// apiLogoutRoute
func apiLogoutRoute(ctx *fiber.Ctx) error {
	s, _ := SessionStore.Get(ctx)
	s.Delete("user_username")
	s.Delete("user_password")
	s.Delete("user_level")
	s.Delete("user_id")
	_ = s.Save()
	return ctx.Redirect("/")
}
