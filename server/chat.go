package server

import (
	"change/database"
	"change/logger"
	"change/tools"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"strings"
)

type ChatWebsocketMessageModel struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// makeChatMessage 生成私聊消息回复
func makeChatMessage(ok bool, message string) []byte {
	var i = ChatWebsocketMessageModel{Message: message}
	if ok {
		i.Status = "ok"
	} else {
		i.Status = "err"
	}
	res, _ := json.Marshal(&i)
	return res
}

// BindChatWebsocket 绑定私聊服务所需要的websocket
func BindChatWebsocket() {
	WebServer.Use("/chat", func(ctx *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(ctx) {
			ctx.Locals("allowed", true)
			s, _ := SessionStore.Get(ctx)
			userID := s.Get("user_id")
			ctx.Locals("user", tools.InterfaceToString(userID))
			return ctx.Next()
		}
		return fiber.ErrUpgradeRequired
	}, middleMustLogin)
	WebServer.Get("/chat/:id", websocket.New(func(c *websocket.Conn) {
		userIDInt := tools.InterfaceToInt(c.Locals("user"))
		id := c.Params("id")
		idInt := tools.StringToInt(id)
		if id == "" || !database.UserHaveUserByID(idInt) {
			return
		}
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				break
			}
			msgString := string(msg)
			logger.ConsoleLogger.Debugln(msgString)
			if strings.HasPrefix(msgString, "chat ") {
				messageSend := strings.TrimPrefix(msgString, "chat ")
				database.MessageCreateChat(userIDInt, idInt, messageSend)
				_ = c.WriteMessage(mt, makeChatMessage(true, "发送成功！"))
			}
		}
	}))
	WebServer.Get("/api/get/chat", func(ctx *fiber.Ctx) error {
		s, _ := SessionStore.Get(ctx)
		userID := s.Get("user_id")
		if userID == nil {
			return ctx.SendString("")
		}
		userIDInt := tools.InterfaceToInt(userID)
		var messages []database.MessageModel
		_ = database.DatabaseEngine.Table(new(database.MessageModel)).Where("type = ?", "chat").Where("to_user = ?", userIDInt).Find(&messages)
		_, _ = database.DatabaseEngine.Table(new(database.MessageModel)).Where("type = ?", "chat").Where("to_user = ?", userIDInt).Delete()
		return ctx.JSON(fiber.Map{"message": messages})
	})
}
