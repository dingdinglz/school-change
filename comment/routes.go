package comment

import (
	"change/database"
	"change/tools"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"xorm.io/xorm"
)

var sessionStore *session.Store

func SetSessionStore(s *session.Store) {
	sessionStore = s
}

// RouteCreate 创建新评论
func RouteCreate(ctx *fiber.Ctx) error {
	var i CommentModel
	_ = ctx.BodyParser(&i)
	s, _ := sessionStore.Get(ctx)
	userID := s.Get("user_id")
	userIDInt := tools.InterfaceToInt(userID)
	i.User = userIDInt
	i.Fullname = database.UserGetRealnameByID(userIDInt)
	CreateNewComment(i)
	return ctx.SendString("ok")
}

// RouteGet 获取所有评论
func RouteGet(ctx *fiber.Ctx) error {
	change := ctx.FormValue("change", "0")
	changeInt := tools.StringToInt(change)
	s, _ := sessionStore.Get(ctx)
	userID := s.Get("user_id")
	userIDInt := tools.InterfaceToInt(userID)
	res := GetAllComments(changeInt, userIDInt)
	if res == nil {
		return ctx.SendString("{\"data\":[]}")
	}
	return ctx.JSON(fiber.Map{"data": res})
}

// RouteLike 为评论点赞
func RouteLike(ctx *fiber.Ctx) error {
	change := ctx.FormValue("change", "0")
	changeInt := tools.StringToInt(change)
	id := ctx.FormValue("id", "")
	s, _ := sessionStore.Get(ctx)
	userID := s.Get("user_id")
	userIDInt := tools.InterfaceToInt(userID)
	dbConnect, _ := xorm.NewEngine("sqlite3", getCommentDB(changeInt))
	_ = dbConnect.Sync2(new(CommentModel))
	var i CommentModel
	_, _ = dbConnect.Table(new(CommentModel)).Where("i_d = ?", id).Get(&i)
	if database.LikeHave(changeInt, id, userIDInt) {
		i.UpvoteCount -= 1
		database.LikeDelete(changeInt, id, userIDInt)
	} else {
		i.UpvoteCount += 1
		database.LikeCreateNew(changeInt, id, userIDInt)
	}
	_, _ = dbConnect.Table(new(CommentModel)).Where("i_d = ?", id).Update(&i)
	_ = dbConnect.Close()
	return ctx.SendString("ok")
}
