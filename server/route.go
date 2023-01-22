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
	WarningMessage := s.Get("message_warning")
	if WarningMessage != nil {
		ResMap["WarningMessage"] = fmt.Sprintf("%s", WarningMessage)
		s.Delete("message_warning")
		_ = s.Save()
	}
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
	WebServer.Get("/user/:id", userRoute)

	changeRoute := WebServer.Group("/change")
	changeRoute.Use(middleMustLogin)
	changeRoute.Get("/new", newChangeRoute)

	apiRoute := WebServer.Group("/api")
	apiRoute.Post("/login", apiLoginRoute)
	apiRoute.Get("/logout", apiLogoutRoute)
	apiRoute.Post("/register", apiRegisterUserRoute)

	apiGetRoute := apiRoute.Group("/get")
	apiGetRoute.Get("/user", apiGetUserInfo)

	apiUploadRoute := apiRoute.Group("/upload")
	apiUploadRoute.Post("/avatar", apiUploadAvatar)

	apiUpdateRoute := apiRoute.Group("/update")
	apiUpdateRoute.Post("/user", apiUpdateUser)

	apiCreateRoute := apiRoute.Group("/create")
	apiCreateRoute.Post("/change", apiCreateChange)

	adminRoute := WebServer.Group("/admin")
	adminRoute.Use(middleAdminRoute)
	adminRoute.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Redirect("/admin/user/manage")
	})
	adminRoute.Get("/user/manage", adminUserManageRoute)
	adminRoute.Get("/user/apply", adminUserApplyRoute)
	adminRoute.Get("/change/subject", adminChangeSubjectRoute)

	adminApiRoute := adminRoute.Group("/api")
	adminApiRoute.Post("/apply/pass", adminApiApplyPass)
	adminApiRoute.Post("/apply/stop", adminApiApplyStop)

	adminApiGetRoute := adminApiRoute.Group("/get")
	adminApiGetRoute.Post("/users", adminApiGetUsers)
	adminApiGetRoute.Get("/applies", adminApiGetApplies)
	adminApiGetRoute.Get("/subjects", adminApiGetSubjects)

	adminApiDeleteRoute := adminApiRoute.Group("/delete")
	adminApiDeleteRoute.Post("/user", adminApiDeleteUserRoute)
	adminApiDeleteRoute.Post("/subject", adminApiDeleteSubject)

	adminApiUpdateRoute := adminApiRoute.Group("/update")
	adminApiUpdateRoute.Post("/user", adminApiUpdateUserRoute)
	adminApiUpdateRoute.Post("/subject", adminApiUpdateSubject)

	adminApiCreateRoute := adminApiRoute.Group("/create")
	adminApiCreateRoute.Post("/user", adminApiCreateUser)
	adminApiCreateRoute.Post("/subject", adminApiCreateSubject)

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

// MakeApiResMapWithData 生成带返回数据的json
func MakeApiResMapWithData(ok bool, message string, data fiber.Map) fiber.Map {
	if ok {
		return fiber.Map{"status": "ok", "message": message, "data": data}
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

// apiLogoutRoute 登出
func apiLogoutRoute(ctx *fiber.Ctx) error {
	s, _ := SessionStore.Get(ctx)
	s.Delete("user_username")
	s.Delete("user_password")
	s.Delete("user_level")
	s.Delete("user_id")
	_ = s.Save()
	return ctx.Redirect("/")
}

// adminUserManageRoute 后台用户管理路由
func adminUserManageRoute(ctx *fiber.Ctx) error {
	return ctx.Render("admin/user", MakeViewMap(ctx), "layout/admin")
}

// adminApiGetUsers 后台取用户列表
func adminApiGetUsers(ctx *fiber.Ctx) error {
	// xorm没有分页的方法，这一段自己写分页逻辑，快把自己整傻了
	page := ctx.FormValue("page")
	pageInt, _ := strconv.Atoi(page)
	if pageInt == 0 {
		return ctx.JSON(MakeApiResMap(false, "page解析错误！"))
	}
	userAll := database.UserGetNumReal()
	var pageAll int = 0
	if userAll%10 == 0 {
		pageAll = userAll / 10
	} else {
		pageAll = userAll/10 + 1
	}
	if pageInt > pageAll {
		return ctx.JSON(MakeApiResMap(false, "page解析错误！"))
	}
	var userAllData []database.UserModel
	_ = database.DatabaseEngine.Table(new(database.UserModel)).Find(&userAllData)
	//logger.ConsoleLogger.Debugln(pageInt*10-10, userAll-1, pageInt*10-1)
	if pageInt*10 >= userAll {
		return ctx.JSON(MakeApiResMapWithData(true, "ok", fiber.Map{"all": userAll, "users": userAllData[(pageInt*10 - 10):(userAll)]}))
	}
	return ctx.JSON(MakeApiResMapWithData(true, "ok", fiber.Map{"all": userAll, "users": userAllData[(pageInt*10 - 10):(pageInt * 10)]}))
}

// adminApiDeleteUserRoute admin删除用户Route
func adminApiDeleteUserRoute(ctx *fiber.Ctx) error {
	id := ctx.FormValue("id", "")
	if id == "" {
		return ctx.JSON(MakeApiResMap(false, "字段为空！"))
	}
	idInt, _ := strconv.Atoi(id)
	if !database.UserHaveUserByID(idInt) {
		return ctx.JSON(MakeApiResMap(false, "用户不存在！"))
	}
	s, _ := SessionStore.Get(ctx)
	myID := s.Get("user_id")
	myIDint, _ := strconv.Atoi(fmt.Sprintf("%s", myID))
	if !database.UserCheckChangeAble(myIDint, idInt) {
		return ctx.JSON(MakeApiResMap(false, "无操作权限！"))
	}
	_, _ = database.DatabaseEngine.Table(new(database.UserModel)).Where("id = ?", idInt).Delete()
	return ctx.JSON(MakeApiResMap(true, "删除成功！"))
}

// adminApiUpdateUserRoute admin更新用户Route
func adminApiUpdateUserRoute(ctx *fiber.Ctx) error {
	id := ctx.FormValue("id", "")
	username := ctx.FormValue("username", "")
	level := ctx.FormValue("level", "")
	grade := ctx.FormValue("grade", "")
	class := ctx.FormValue("class", "")
	realname := ctx.FormValue("realname", "")
	s, _ := SessionStore.Get(ctx)
	myID := s.Get("user_id")
	if !database.UserCheckChangeAble(tools.InterfaceToInt(myID), tools.StringToInt(id)) {
		return ctx.JSON(MakeApiResMap(false, "无操作权限！"))
	}
	if tools.StringToInt(level) > database.UserGetLevelByID(tools.InterfaceToInt(myID)) {
		return ctx.JSON(MakeApiResMap(false, "不能操作到比自身高的level！"))
	}
	if id == "" || username == "" || level == "" || grade == "" || class == "" || realname == "" {
		return ctx.JSON(MakeApiResMap(false, "存在字段为空！"))
	}
	_, _ = database.DatabaseEngine.Table(new(database.UserModel)).Where("id = ?", tools.StringToInt(id)).Update(&database.UserModel{
		Username: username,
		Level:    tools.StringToInt(level),
		Grade:    tools.StringToInt(grade),
		Class:    tools.StringToInt(class),
		Realname: realname,
	})
	return ctx.JSON(MakeApiResMap(true, "更新成功！"))
}

// adminApiCreateUser 创建新用户
func adminApiCreateUser(ctx *fiber.Ctx) error {
	username := ctx.FormValue("username", "")
	password := ctx.FormValue("password", "")
	level := ctx.FormValue("level", "")
	grade := ctx.FormValue("grade", "")
	class := ctx.FormValue("class", "")
	realname := ctx.FormValue("realname", "")
	ip := ctx.IP()
	if username == "" || password == "" || level == "" || grade == "" || class == "" || realname == "" {
		return ctx.JSON(MakeApiResMap(false, "存在字段为空！"))
	}
	_, e := database.UserCheckCreateAble(username, password, tools.StringToInt(level), tools.StringToInt(class), tools.StringToInt(grade))
	if e != nil {
		return ctx.JSON(MakeApiResMap(false, e.Error()))
	}
	s, _ := SessionStore.Get(ctx)
	myID := s.Get("user_id")
	if database.UserGetLevelByID(tools.InterfaceToInt(myID)) < tools.StringToInt(level) {
		return ctx.JSON(MakeApiResMap(false, "权限不足！"))
	}
	err := database.UserCreateNew(username, password, tools.StringToInt(level), ip, tools.StringToInt(grade), tools.StringToInt(class), realname)
	if err != nil {
		return ctx.JSON(MakeApiResMap(false, "创建用户失败！"+err.Error()))
	}
	return ctx.JSON(MakeApiResMap(true, "创建成功！"))
}

// apiRegisterUserRoute 用户注册api
func apiRegisterUserRoute(ctx *fiber.Ctx) error {
	username := ctx.FormValue("username", "")
	password := ctx.FormValue("password", "")
	grade := ctx.FormValue("grade", "")
	class := ctx.FormValue("class", "")
	realname := ctx.FormValue("realname", "")
	ip := ctx.IP()
	if username == "" || password == "" || grade == "" || class == "" || realname == "" {
		return ctx.JSON(MakeApiResMap(false, "存在字段为空！"))
	}
	_, e := database.UserCheckCreateAble(username, password, 1, tools.StringToInt(class), tools.StringToInt(grade))
	if e != nil {
		return ctx.JSON(MakeApiResMap(false, e.Error()))
	}
	if database.UserApplyHaveIP(ip) {
		return ctx.JSON(MakeApiResMap(false, "该ip存在账号申请！"))
	}
	database.UserApplyCreate(username, password, ip, tools.StringToInt(grade), tools.StringToInt(class), realname)
	return ctx.JSON(MakeApiResMap(true, "已申请！请等待审核！"))
}

// adminApiGetApplies 取所有apply
func adminApiGetApplies(ctx *fiber.Ctx) error {
	var appliesData []database.ApplyModel
	_ = database.DatabaseEngine.Table(new(database.ApplyModel)).Find(&appliesData)
	return ctx.JSON(MakeApiResMapWithData(true, "获取成功！", fiber.Map{"applies": appliesData}))
}

// adminUserApplyRoute 用户申请route
func adminUserApplyRoute(ctx *fiber.Ctx) error {
	return ctx.Render("admin/apply", MakeViewMap(ctx), "layout/admin")
}

// adminApiApplyPass 通过用户申请
func adminApiApplyPass(ctx *fiber.Ctx) error {
	ip := ctx.FormValue("ip", "")
	if ip == "" {
		return ctx.JSON(MakeApiResMap(false, "存在字段为空！"))
	}
	e := database.UserApplyPass(ip)
	if e != nil {
		return ctx.JSON(MakeApiResMap(false, e.Error()))
	}
	return ctx.JSON(MakeApiResMap(true, "通过成功！"))
}

// adminApiApplyStop 拒绝用户申请
func adminApiApplyStop(ctx *fiber.Ctx) error {
	ip := ctx.FormValue("ip", "")
	if ip == "" {
		return ctx.JSON(MakeApiResMap(false, "存在字段为空！"))
	}
	database.UserApplyStop(ip)
	return ctx.JSON(MakeApiResMap(true, "拒绝成功！"))
}

// userRoute 用户资料route
func userRoute(ctx *fiber.Ctx) error {
	viewMap := MakeViewMap(ctx)
	viewMap["View_ID"] = tools.StringToInt(ctx.Params("id", "0"))
	return ctx.Render("user", viewMap, "layout/main")
}

// apiGetUserInfo 获取用户信息
func apiGetUserInfo(ctx *fiber.Ctx) error {
	id := ctx.Query("id", "")
	if id == "" {
		return ctx.JSON(MakeApiResMap(false, "id为空！"))
	}
	if !database.UserHaveUserByID(tools.StringToInt(id)) {
		return ctx.JSON(MakeApiResMap(false, "用户不存在！"))
	}
	var u database.UserModel
	_, _ = database.DatabaseEngine.Table(new(database.UserModel)).Where("id = ?", tools.StringToInt(id)).Get(&u)
	u.Ip = ""
	u.Password = ""
	return ctx.JSON(MakeApiResMapWithData(true, "获取成功！", fiber.Map{"user": u}))
}

// apiUploadAvatar 上传新头像
func apiUploadAvatar(ctx *fiber.Ctx) error {
	s, _ := SessionStore.Get(ctx)
	UserUsername := s.Get("user_username")
	UserPassword := s.Get("user_password")
	UserLevel := s.Get("user_level")
	UserID := s.Get("user_id")
	if UserUsername == nil || UserPassword == nil || UserLevel == nil {
		return ctx.JSON(MakeApiResMap(false, "请先登录！"))
	} else {
		UserLevel_i, _ := strconv.Atoi(fmt.Sprintf("%s", UserLevel))
		UserID_i, _ := strconv.Atoi(fmt.Sprintf("%s", UserID))
		// logger.ConsoleLogger.Debugln(fmt.Sprintf("%s", UserUsername), fmt.Sprintf("%s", UserPassword), UserLevel_i)
		if database.UserCheck(fmt.Sprintf("%s", UserUsername), fmt.Sprintf("%s", UserPassword), UserLevel_i, UserID_i) {
			rootPath, _ := os.Getwd()
			f, _ := ctx.FormFile("file")
			_ = ctx.SaveFile(f, filepath.Join(rootPath, "data", "avatar", fmt.Sprintf("%s", UserUsername)+".png"))
			return ctx.JSON(MakeApiResMap(true, "上传成功！"))
		}
	}
	return ctx.JSON(MakeApiResMap(false, "请重新登录！"))
}

// apiUpdateUser 更新用户资料
func apiUpdateUser(ctx *fiber.Ctx) error {
	s, _ := SessionStore.Get(ctx)
	UserUsername := s.Get("user_username")
	UserPassword := s.Get("user_password")
	UserLevel := s.Get("user_level")
	UserID := s.Get("user_id")
	if UserUsername == nil || UserPassword == nil || UserLevel == nil {
		return ctx.JSON(MakeApiResMap(false, "请先登录！"))
	} else {
		UserLevel_i, _ := strconv.Atoi(fmt.Sprintf("%s", UserLevel))
		UserID_i, _ := strconv.Atoi(fmt.Sprintf("%s", UserID))
		// logger.ConsoleLogger.Debugln(fmt.Sprintf("%s", UserUsername), fmt.Sprintf("%s", UserPassword), UserLevel_i)
		if database.UserCheck(fmt.Sprintf("%s", UserUsername), fmt.Sprintf("%s", UserPassword), UserLevel_i, UserID_i) {
			old := ctx.FormValue("oldpassword", "")
			newP := ctx.FormValue("newpassword", "")
			newPagain := ctx.FormValue("newpasswordagain", "")
			if old == "" || newP == "" || newPagain == "" {
				return ctx.JSON(MakeApiResMap(false, "存在字段为空！"))
			}
			if tools.MD5(old) != fmt.Sprintf("%s", UserPassword) {
				return ctx.JSON(MakeApiResMap(false, "老密码错误！"))
			}
			if newP != newPagain {
				return ctx.JSON(MakeApiResMap(false, "新密码两次输入不一致！"))
			}
			_, _ = database.DatabaseEngine.Table(new(database.UserModel)).Where("id = ?", UserID_i).Update(&database.UserModel{Password: tools.MD5(newP)})
			return ctx.JSON(MakeApiResMap(true, "更新成功！"))
		}
	}
	return ctx.JSON(MakeApiResMap(false, "请重新登录！"))
}

// adminApiGetSubjects 取subject列表
func adminApiGetSubjects(ctx *fiber.Ctx) error {
	var subjects []database.SubjectModel
	_ = database.DatabaseEngine.Table(new(database.SubjectModel)).Find(&subjects)
	return ctx.JSON(MakeApiResMapWithData(true, "获取成功！", fiber.Map{"subjects": subjects}))
}

// adminApiCreateSubject 创建subject
func adminApiCreateSubject(ctx *fiber.Ctx) error {
	name := ctx.FormValue("name", "")
	description := ctx.FormValue("description", "")
	if name == "" || description == "" {
		return ctx.JSON(MakeApiResMap(false, "存在字段为空！"))
	}
	database.SubjectCreateNew(name, description)
	return ctx.JSON(MakeApiResMap(true, "创建成功！"))
}

// adminApiDeleteSubject 删除subject
func adminApiDeleteSubject(ctx *fiber.Ctx) error {
	id := ctx.FormValue("id", "")
	if id == "" {
		return ctx.JSON(MakeApiResMap(false, "存在字段为空！"))
	}
	_, _ = database.DatabaseEngine.Table(new(database.SubjectModel)).Where("id = ?", tools.StringToInt(id)).Delete()
	return ctx.JSON(MakeApiResMap(true, "删除成功！"))
}

// adminApiUpdateSubject 更新subject
func adminApiUpdateSubject(ctx *fiber.Ctx) error {
	id := ctx.FormValue("id", "")
	description := ctx.FormValue("description", "")
	if id == "" || description == "" {
		return ctx.JSON(MakeApiResMap(false, "存在字段为空！"))
	}
	_, _ = database.DatabaseEngine.Table(new(database.SubjectModel)).Where("id = ?", tools.StringToInt(id)).Update(&database.SubjectModel{Description: description})
	return ctx.JSON(MakeApiResMap(true, "更新成功！"))
}

// adminChangeSubjectRoute 后台管理subject路由
func adminChangeSubjectRoute(ctx *fiber.Ctx) error {
	return ctx.Render("admin/subject", MakeViewMap(ctx), "layout/admin")
}

// newChangeRoute 新建change
func newChangeRoute(ctx *fiber.Ctx) error {
	viewMap := MakeViewMap(ctx)
	var subjects []database.SubjectModel
	_ = database.DatabaseEngine.Table(new(database.SubjectModel)).Find(&subjects)
	viewMap["Subjects"] = subjects
	return ctx.Render("newchange", viewMap, "layout/main")
}

// apiCreateChange 创建change
func apiCreateChange(ctx *fiber.Ctx) error {
	title := ctx.FormValue("title", "")
	description := ctx.FormValue("description", "")
	subject := ctx.FormValue("subject", "")
	want := ctx.FormValue("want", "")
	if title == "" || description == "" || subject == "" || want == "" {
		return ctx.JSON(MakeApiResMap(false, "存在字段为空！"))
	}
	s, _ := SessionStore.Get(ctx)
	userID := s.Get("user_id")
	userIDStr := tools.InterfaceToString(userID)
	if !database.ChangeCheckTime(tools.StringToInt(userIDStr)) {
		return ctx.JSON(MakeApiResMap(false, "距离上次创建不足半小时！"))
	}
	database.ChangeCreateNew(title, description, tools.StringToInt(subject), tools.StringToInt(userIDStr), want)
	return ctx.JSON(MakeApiResMap(true, "发布成功！"))
}
