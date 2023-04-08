package server

import (
	"change/comment"
	"change/config"
	"change/database"
	"change/tools"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
			ResMap["MessageNum"], _ = database.DatabaseEngine.Table(new(database.MessageModel)).Where("to_user = ?", UserID_i).Count()
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
	WebServer.Get("/money/high", moneyHighRoute)
	WebServer.Get("/money/low", moneyLowRoute)
	WebServer.Get("/avatar/:username", avatarRoute)
	WebServer.Get("/user/:id", userRoute)
	WebServer.Get("/picture/change/:change/:which", changePictureRoute)
	WebServer.Get("/subject", subjectRoute)
	WebServer.Get("/subject/id/:id", subjectAllRoute)
	WebServer.Use("/webchat", middleMustLogin)
	WebServer.Get("/webchat/:id", chatRoute)
	WebServer.Use("/info", middleMustLogin)
	WebServer.Get("/info", messageRoute)
	WebServer.Get("/about", aboutRoute)
	WebServer.Get("/status", monitor.New(monitor.Config{Title: "monitor"}))
	WebServer.Get("/search", searchRoute)

	changeRoute := WebServer.Group("/change")
	changeRoute.Use(middleMustLogin)
	changeRoute.Get("/new", newChangeRoute)
	changeRoute.Get("/my", myChangeRoute)
	changeRoute.Get("/id/:id", changeSeeRoute)
	changeRoute.Post("/uploadImage", apiChangeUploadImage)
	changeRoute.Post("/deleteImage", apiChangeDeleteImage)
	changeRoute.Post("/delete", apiChangeDelete)
	changeRoute.Post("/report", apiChangeReportRoute)
	changeRoute.Get("/cleanMessages", apiCleanMessages)

	changeApiRoute := changeRoute.Group("/api")

	changeApiUpdateRoute := changeApiRoute.Group("/update")
	changeApiUpdateRoute.Post("/changeState", apiUpdateChangeState)

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

	commentRoute := apiRoute.Group("/comment")
	commentRoute.Use(middleMustLogin)
	commentRoute.Post("/post", comment.RouteCreate)
	commentRoute.Post("/get", comment.RouteGet)
	commentRoute.Post("/like", comment.RouteLike)

	adminRoute := WebServer.Group("/admin")
	adminRoute.Use(middleAdminRoute)
	adminRoute.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Redirect("/admin/user/manage")
	})
	adminRoute.Get("/user/manage", adminUserManageRoute)
	adminRoute.Get("/user/apply", adminUserApplyRoute)
	adminRoute.Get("/change/subject", adminChangeSubjectRoute)
	adminRoute.Get("/change/report", adminReportRoute)
	adminRoute.Get("/status", statusRoute)

	adminApiRoute := adminRoute.Group("/api")
	adminApiRoute.Post("/apply/pass", adminApiApplyPass)
	adminApiRoute.Post("/apply/stop", adminApiApplyStop)

	adminApiGetRoute := adminApiRoute.Group("/get")
	adminApiGetRoute.Post("/users", adminApiGetUsers)
	adminApiGetRoute.Get("/applies", adminApiGetApplies)
	adminApiGetRoute.Get("/subjects", adminApiGetSubjects)
	adminApiGetRoute.Get("/reports", apiAdminGetReports)

	adminApiDeleteRoute := adminApiRoute.Group("/delete")
	adminApiDeleteRoute.Post("/user", adminApiDeleteUserRoute)
	adminApiDeleteRoute.Post("/subject", adminApiDeleteSubject)
	adminApiDeleteRoute.Post("/report", adminDeleteReport)

	adminApiUpdateRoute := adminApiRoute.Group("/update")
	adminApiUpdateRoute.Post("/user", adminApiUpdateUserRoute)
	adminApiUpdateRoute.Post("/subject", adminApiUpdateSubject)

	adminApiCreateRoute := adminApiRoute.Group("/create")
	adminApiCreateRoute.Post("/user", adminApiCreateUser)
	adminApiCreateRoute.Post("/subject", adminApiCreateSubject)

	adminApiPassRoute := adminApiRoute.Group("/pass")
	adminApiPassRoute.Post("/report", adminReportPass)
}

// indexRoute 主页路由
func indexRoute(ctx *fiber.Ctx) error {
	nowpage := ctx.Query("page", "1")
	nowpageInt := tools.StringToInt(nowpage)
	viewMap := MakeViewMap(ctx)
	var allChanges, Changes []database.ChangeModel
	var pages int = 0
	_ = database.DatabaseEngine.Table(new(database.ChangeModel)).Where("state = ?", 1).Desc("time").Find(&Changes)
	if len(Changes) == 0 {
		return ctx.Render("index", viewMap, "layout/main")
	}
	pages = len(Changes) / 12
	if len(Changes)%12 != 0 {
		pages += 1
	}
	if nowpageInt <= 0 || nowpageInt > pages {
		nowpageInt = 1
	}
	showChanges := [][]database.ChangeModel{}
	paginations := []string{}
	for i := 1; i <= pages; i++ {
		paginations = append(paginations, strconv.Itoa(i))
	}

	if pages != nowpageInt {
		allChanges = Changes[(nowpageInt*12 - 12):(nowpageInt * 12)]
	} else {
		allChanges = Changes[(nowpageInt*12 - 12):]
	}
	for i := 1; i <= len(allChanges); i += 4 {
		if i+4 > len(allChanges) {
			showChanges = append(showChanges, allChanges[(i-1):])
		} else {
			showChanges = append(showChanges, allChanges[(i-1):i+3])
		}
	}
	viewMap["Changes"] = showChanges
	viewMap["Page_map"] = paginations
	viewMap["Page_now"] = nowpage
	viewMap["Page_all"] = strconv.Itoa(pages)
	viewMap["Page_next"] = strconv.Itoa(nowpageInt + 1)
	viewMap["Page_present"] = strconv.Itoa(nowpageInt - 1)
	return ctx.Render("index", viewMap, "layout/main")
}

// moneyHighRoute 按金额排序由高到低
func moneyHighRoute(ctx *fiber.Ctx) error {
	nowpage := ctx.Query("page", "1")
	nowpageInt := tools.StringToInt(nowpage)
	viewMap := MakeViewMap(ctx)
	var allChanges, Changes []database.ChangeModel
	var pages int = 0
	_ = database.DatabaseEngine.Table(new(database.ChangeModel)).Where("state = ?", 1).Desc("money").Desc("time").Find(&Changes)
	if len(Changes) == 0 {
		return ctx.Render("index", viewMap, "layout/main")
	}
	pages = len(Changes) / 12
	if len(Changes)%12 != 0 {
		pages += 1
	}
	if nowpageInt <= 0 || nowpageInt > pages {
		nowpageInt = 1
	}
	showChanges := [][]database.ChangeModel{}
	paginations := []string{}
	for i := 1; i <= pages; i++ {
		paginations = append(paginations, strconv.Itoa(i))
	}

	if pages != nowpageInt {
		allChanges = Changes[(nowpageInt*12 - 12):(nowpageInt * 12)]
	} else {
		allChanges = Changes[(nowpageInt*12 - 12):]
	}
	for i := 1; i <= len(allChanges); i += 4 {
		if i+4 > len(allChanges) {
			showChanges = append(showChanges, allChanges[(i-1):])
		} else {
			showChanges = append(showChanges, allChanges[(i-1):i+3])
		}
	}
	viewMap["Changes"] = showChanges
	viewMap["Page_map"] = paginations
	viewMap["Page_now"] = nowpage
	viewMap["Page_all"] = strconv.Itoa(pages)
	viewMap["Page_next"] = strconv.Itoa(nowpageInt + 1)
	viewMap["Page_present"] = strconv.Itoa(nowpageInt - 1)
	return ctx.Render("index", viewMap, "layout/main")
}

// moneyLowRoute 按金额排序由低到高
func moneyLowRoute(ctx *fiber.Ctx) error {
	nowpage := ctx.Query("page", "1")
	nowpageInt := tools.StringToInt(nowpage)
	viewMap := MakeViewMap(ctx)
	var allChanges, Changes []database.ChangeModel
	var pages int = 0
	_ = database.DatabaseEngine.Table(new(database.ChangeModel)).Where("state = ?", 1).Asc("money").Desc("time").Find(&Changes)
	if len(Changes) == 0 {
		return ctx.Render("index", viewMap, "layout/main")
	}
	pages = len(Changes) / 12
	if len(Changes)%12 != 0 {
		pages += 1
	}
	if nowpageInt <= 0 || nowpageInt > pages {
		nowpageInt = 1
	}
	showChanges := [][]database.ChangeModel{}
	paginations := []string{}
	for i := 1; i <= pages; i++ {
		paginations = append(paginations, strconv.Itoa(i))
	}

	if pages != nowpageInt {
		allChanges = Changes[(nowpageInt*12 - 12):(nowpageInt * 12)]
	} else {
		allChanges = Changes[(nowpageInt*12 - 12):]
	}
	for i := 1; i <= len(allChanges); i += 4 {
		if i+4 > len(allChanges) {
			showChanges = append(showChanges, allChanges[(i-1):])
		} else {
			showChanges = append(showChanges, allChanges[(i-1):i+3])
		}
	}
	viewMap["Changes"] = showChanges
	viewMap["Page_map"] = paginations
	viewMap["Page_now"] = nowpage
	viewMap["Page_all"] = strconv.Itoa(pages)
	viewMap["Page_next"] = strconv.Itoa(nowpageInt + 1)
	viewMap["Page_present"] = strconv.Itoa(nowpageInt - 1)
	return ctx.Render("index", viewMap, "layout/main")
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
	database.MessageCreateToAdmins("apply", "用户注册申请："+realname)
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
	money := ctx.FormValue("money", "")
	if title == "" || description == "" || subject == "" || want == "" || money == "" {
		return ctx.JSON(MakeApiResMap(false, "存在字段为空！"))
	}
	s, _ := SessionStore.Get(ctx)
	userID := s.Get("user_id")
	userIDStr := tools.InterfaceToString(userID)
	if !database.ChangeCheckTime(tools.StringToInt(userIDStr)) {
		return ctx.JSON(MakeApiResMap(false, "距离上次创建不足"+strconv.Itoa(config.ConfigData.Limit.Change)+"分钟！"))
	}
	database.ChangeCreateNew(title, description, tools.StringToInt(subject), tools.StringToInt(userIDStr), want, tools.StringToInt(money))
	return ctx.JSON(MakeApiResMap(true, "发布成功！"))
}

// myChangeRoute 我的交换路由
func myChangeRoute(ctx *fiber.Ctx) error {
	s, _ := SessionStore.Get(ctx)
	userID := s.Get("user_id")
	userIDStr := tools.InterfaceToString(userID)
	viewMap := MakeViewMap(ctx)
	var myChanges []database.ChangeModel
	_ = database.DatabaseEngine.Table(new(database.ChangeModel)).Where("user = ?", tools.StringToInt(userIDStr)).Desc("time").Find(&myChanges)
	viewMap["Changes"] = myChanges
	return ctx.Render("showchange", viewMap, "layout/main")
}

// changePictureRoute 交换图片路由
func changePictureRoute(ctx *fiber.Ctx) error {
	change := ctx.Params("change", "")
	which := ctx.Params("which", "")
	rootPath, _ := os.Getwd()
	if !tools.IsFileExist(filepath.Join(rootPath, "data", "change_pic", change, which+".png")) {
		return ctx.SendFile(filepath.Join(rootPath, "web", "picture", "change.png"))
	}
	return ctx.SendFile(filepath.Join(rootPath, "data", "change_pic", change, which+".png"))
}

// changeSeeRoute 查看详情route
func changeSeeRoute(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	idInt := tools.StringToInt(id)
	var c database.ChangeModel
	_, _ = database.DatabaseEngine.Table(new(database.ChangeModel)).Where("id = ?", idInt).Get(&c)
	viewMap := MakeViewMap(ctx)
	viewMap["Change"] = c
	rootPath, _ := os.Getwd()
	imagePath := filepath.Join(rootPath, "data", "change_pic", id)
	if !tools.IsFileExist(imagePath) {
		viewMap["Images"] = []fiber.Map{}
	} else {
		addImages := []fiber.Map{}
		_ = filepath.Walk(imagePath, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			addImages = append(addImages, fiber.Map{"Id": id, "Name": strings.ReplaceAll(filepath.Base(path), ".png", "")})
			return nil
		})
		viewMap["Images"] = addImages
	}
	return ctx.Render("seechange", viewMap, "layout/main")
}

// apiUpdateChangeState 更新状态
func apiUpdateChangeState(ctx *fiber.Ctx) error {
	id := ctx.FormValue("id", "")
	state := ctx.FormValue("state", "")
	if id == "" || state == "" {
		return ctx.JSON(MakeApiResMap(false, "存在字段为空！"))
	}
	if !database.ChangeHaveByID(tools.StringToInt(id)) {
		return ctx.JSON(MakeApiResMap(false, "该交换不存在！"))
	}
	s, _ := SessionStore.Get(ctx)
	userID := s.Get("user_id")
	userIDInt := tools.InterfaceToInt(userID)
	var c database.ChangeModel
	_, _ = database.DatabaseEngine.Table(new(database.ChangeModel)).Where("id = ?", tools.StringToInt(id)).Get(&c)
	if c.User != userIDInt {
		return ctx.JSON(MakeApiResMap(false, "你无权操作一个不属于你的交换！"))
	}
	_, _ = database.DatabaseEngine.Table(new(database.ChangeModel)).Where("id = ?", tools.StringToInt(id)).Update(&database.ChangeModel{State: tools.StringToInt(state)})
	return ctx.JSON(MakeApiResMap(true, "操作成功！"))
}

// apiChangeUploadImage 上传change照片
func apiChangeUploadImage(ctx *fiber.Ctx) error {
	id := ctx.Query("id", "")
	if id == "" || !database.ChangeHaveByID(tools.StringToInt(id)) {
		return ctx.JSON(MakeApiResMap(false, "交换不存在！"))
	}
	s, _ := SessionStore.Get(ctx)
	userID := s.Get("user_id")
	userIDInt := tools.InterfaceToInt(userID)
	if !database.ChangeCheckUser(tools.StringToInt(id), userIDInt) {
		return ctx.JSON(MakeApiResMap(false, "你无权操作一个不属于你的交换！"))

	}
	rootPath, _ := os.Getwd()
	if !tools.IsFileExist(filepath.Join(rootPath, "data", "change_pic", id)) {
		_ = os.Mkdir(filepath.Join(rootPath, "data", "change_pic", id), 0777)
	}
	saveName := tools.GetNextPicName(filepath.Join(rootPath, "data", "change_pic", id))
	f, _ := ctx.FormFile("file")
	_ = ctx.SaveFile(f, filepath.Join(rootPath, "data", "change_pic", id, saveName))
	return ctx.JSON(MakeApiResMap(true, "上传成功！"))
}

func apiChangeDeleteImage(ctx *fiber.Ctx) error {
	id := ctx.FormValue("id", "")
	name := ctx.FormValue("name", "")
	if id == "" || name == "" {
		return ctx.JSON(MakeApiResMap(false, "存在字段为空！"))
	}
	if !database.ChangeHaveByID(tools.StringToInt(id)) {
		return ctx.JSON(MakeApiResMap(false, "交换不存在！"))
	}
	s, _ := SessionStore.Get(ctx)
	userID := s.Get("user_id")
	userIDInt := tools.InterfaceToInt(userID)
	if !database.ChangeCheckUser(tools.StringToInt(id), userIDInt) {
		return ctx.JSON(MakeApiResMap(false, "你无权操作一个不属于你的交换！"))
	}
	rootPath, _ := os.Getwd()
	if !tools.IsFileExist(filepath.Join(rootPath, "data", "change_pic", id, name+".png")) {
		return ctx.JSON(MakeApiResMap(false, "图片不存在！"))
	}
	_ = os.Remove(filepath.Join(rootPath, "data", "change_pic", id, name+".png"))
	return ctx.JSON(MakeApiResMap(true, "删除成功！"))
}

// apiChangeDelete 删除交换
func apiChangeDelete(ctx *fiber.Ctx) error {
	id := ctx.FormValue("id", "")
	if id == "" {
		return ctx.JSON(MakeApiResMap(false, "存在字段为空！"))
	}
	s, _ := SessionStore.Get(ctx)
	userID := s.Get("user_id")
	userIDInt := tools.InterfaceToInt(userID)
	if !database.ChangeCheckUser(tools.StringToInt(id), userIDInt) {
		return ctx.JSON(MakeApiResMap(false, "你无权操作一个不属于你的交换！"))
	}
	if !database.ChangeHaveByID(tools.StringToInt(id)) {
		return ctx.JSON(MakeApiResMap(false, "交换不存在！"))
	}
	_, _ = database.DatabaseEngine.Table(new(database.ChangeModel)).Where("id = ?", tools.StringToInt(id)).Delete()
	return ctx.JSON(MakeApiResMap(true, "删除成功！"))
}

// subjectRoute 分类路由
func subjectRoute(ctx *fiber.Ctx) error {
	viewMap := MakeViewMap(ctx)
	var subjects []database.SubjectModel
	_ = database.DatabaseEngine.Table(new(database.SubjectModel)).Find(&subjects)
	viewMap["Subjects"] = subjects
	return ctx.Render("subject", viewMap, "layout/main")
}

func subjectAllRoute(ctx *fiber.Ctx) error {
	nowpage := ctx.Query("page", "1")
	nowpageInt := tools.StringToInt(nowpage)
	viewMap := MakeViewMap(ctx)
	var allChanges, Changes []database.ChangeModel
	var pages int = 0
	_ = database.DatabaseEngine.Table(new(database.ChangeModel)).Where("subject = ?", tools.StringToInt(ctx.Params("id"))).Where("state = ?", 1).Desc("time").Find(&Changes)
	pages = len(Changes) / 12
	if len(Changes)%12 != 0 {
		pages += 1
	}
	if nowpageInt <= 0 || nowpageInt > pages {
		nowpageInt = 1
	}
	showChanges := [][]database.ChangeModel{}
	paginations := []string{}
	for i := 1; i <= pages; i++ {
		paginations = append(paginations, strconv.Itoa(i))
	}

	if pages != nowpageInt {
		allChanges = Changes[(nowpageInt*12 - 12):(nowpageInt * 12)]
	} else {
		allChanges = Changes[(nowpageInt*12 - 12):]
	}
	for i := 1; i <= len(allChanges); i += 4 {
		if i+4 > len(allChanges) {
			showChanges = append(showChanges, allChanges[(i-1):])
		} else {
			showChanges = append(showChanges, allChanges[(i-1):i+3])
		}
	}
	viewMap["Changes"] = showChanges
	viewMap["Page_map"] = paginations
	viewMap["Page_now"] = nowpage
	viewMap["Page_all"] = strconv.Itoa(pages)
	viewMap["Page_next"] = strconv.Itoa(nowpageInt + 1)
	viewMap["Page_present"] = strconv.Itoa(nowpageInt - 1)
	return ctx.Render("index", viewMap, "layout/main")
}

// chatRoute 聊天界面路由
func chatRoute(ctx *fiber.Ctx) error {
	viewMap := MakeViewMap(ctx)
	viewMap["ChatID"] = tools.StringToInt(ctx.Params("id"))
	return ctx.Render("chat", viewMap)
}

// messageRoute 消息路由
func messageRoute(ctx *fiber.Ctx) error {
	s, _ := SessionStore.Get(ctx)
	userID := s.Get("user_id")
	userIDInt := tools.InterfaceToInt(userID)
	viewMap := MakeViewMap(ctx)
	var messages []database.MessageModel
	_ = database.DatabaseEngine.Table(new(database.MessageModel)).Where("to_user = ?", userIDInt).Desc("time").Find(&messages)
	viewMap["Messages"] = messages
	return ctx.Render("message", viewMap, "layout/main")
}

// apiChangeReportRoute 举报交换api
func apiChangeReportRoute(ctx *fiber.Ctx) error {
	change := ctx.FormValue("change", "")
	message := ctx.FormValue("message", "")
	if change == "" || message == "" {
		return ctx.JSON(MakeApiResMap(false, "存在字段为空！"))
	}
	s, _ := SessionStore.Get(ctx)
	userID := s.Get("user_id")
	userIDInt := tools.InterfaceToInt(userID)
	database.ReportCreateNew(tools.StringToInt(change), message, userIDInt)
	database.MessageCreateToAdmins("report", "新举报："+message)
	return ctx.JSON(MakeApiResMap(true, "举报成功！"))
}

// apiAdminGetReports admin获取举报
func apiAdminGetReports(ctx *fiber.Ctx) error {
	var reports []database.ReportModel
	_ = database.DatabaseEngine.Table(new(database.ReportModel)).Find(&reports)
	return ctx.JSON(MakeApiResMapWithData(true, "获取成功！", fiber.Map{"reports": reports}))
}

// adminReportRoute admin举报路由
func adminReportRoute(ctx *fiber.Ctx) error {
	return ctx.Render("admin/report", MakeViewMap(ctx), "layout/admin")
}

// adminReportPass 通过举报
func adminReportPass(ctx *fiber.Ctx) error {
	change := ctx.FormValue("change", "")
	message := ctx.FormValue("message", "")
	user := ctx.FormValue("user", "")
	if change == "" || message == "" || user == "" {
		return ctx.JSON(MakeApiResMap(false, "存在字段为空！"))
	}
	changeID := tools.StringToInt(change)
	if !database.ChangeHaveByID(changeID) {
		return ctx.JSON(MakeApiResMap(false, "该交换不存在！"))
	}
	var c database.ChangeModel
	_, _ = database.DatabaseEngine.Table(new(database.ChangeModel)).Where("id = ?", changeID).Get(&c)
	s, _ := SessionStore.Get(ctx)
	userID := s.Get("user_id")
	userIDInt := tools.InterfaceToInt(userID)
	if database.UserGetLevelByID(c.User) >= 2 && database.UserGetLevelByID(userIDInt) < database.UserGetLevelByID(c.User) {
		if database.UserGetLevelByID(userIDInt) != 3 {
			return ctx.JSON(MakeApiResMap(false, "权限不足！"))
		}
	}
	_, _ = database.DatabaseEngine.Table(new(database.MessageModel)).Insert(&database.MessageModel{
		Time:     time.Now(),
		Type:     "punish",
		FromUser: 0,
		ToUser:   c.User,
		Message:  message,
	})
	_, _ = database.DatabaseEngine.Table(new(database.MessageModel)).Insert(&database.MessageModel{
		Time:     time.Now(),
		Type:     "thank",
		FromUser: 0,
		ToUser:   tools.StringToInt(user),
		Message:  "您对于交换id为" + change + "的举报成功！感谢反馈！",
	})
	_, _ = database.DatabaseEngine.Table(new(database.ReportModel)).Where("change = ?", changeID).Delete()
	_, _ = database.DatabaseEngine.Table(new(database.ChangeModel)).Where("id = ?", changeID).Delete()
	return ctx.JSON(MakeApiResMap(true, "已通过举报！"))
}

// adminDeleteReport 取消举报
func adminDeleteReport(ctx *fiber.Ctx) error {
	change := ctx.FormValue("change", "")
	message := ctx.FormValue("message", "")
	if change == "" || message == "" {
		return ctx.JSON(MakeApiResMap(false, "存在字段为空！"))
	}
	_, _ = database.DatabaseEngine.Table(new(database.ReportModel)).Where("change = ?", tools.StringToInt(change)).Where("message = ?", message).Delete()
	return ctx.JSON(MakeApiResMap(true, "已取消！"))
}

// apiCleanMessages 清除消息列表
func apiCleanMessages(ctx *fiber.Ctx) error {
	s, _ := SessionStore.Get(ctx)
	userID := s.Get("user_id")
	userIDInt := tools.InterfaceToInt(userID)
	_, _ = database.DatabaseEngine.Table(new(database.MessageModel)).Where("to_user = ?", userIDInt).Delete()
	return ctx.JSON(MakeApiResMap(true, "删除成功！"))
}

// aboutRoute 关于路由
func aboutRoute(ctx *fiber.Ctx) error {
	return ctx.Render("about", MakeViewMap(ctx), "layout/main")
}

// statusRoute 状态路由
func statusRoute(ctx *fiber.Ctx) error {
	return ctx.Render("admin/status", MakeViewMap(ctx), "layout/admin")
}

// searchRoute 搜索路由
func searchRoute(ctx *fiber.Ctx) error {
	word := ctx.Query("word", "")
	if word == "" {
		s, _ := SessionStore.Get(ctx)
		s.Set("message_warning", "登录信息有误！请重新登录！")
		_ = s.Save()
		return ctx.Redirect("/")
	}
	viewMap := MakeViewMap(ctx)
	var allChanges []database.ChangeModel
	_ = database.DatabaseEngine.Table(new(database.ChangeModel)).Where("description like ?", "%"+word+"%").Where("state = ?", 1).Desc("time").Find(&allChanges)
	showChanges := [][]database.ChangeModel{}
	for i := 1; i <= len(allChanges); i += 4 {
		if i+4 > len(allChanges) {
			showChanges = append(showChanges, allChanges[(i-1):])
		} else {
			showChanges = append(showChanges, allChanges[(i-1):i+3])
		}
	}
	viewMap["Changes"] = showChanges
	return ctx.Render("index", viewMap, "layout/main")

}
