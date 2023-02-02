package database

import "time"

// MessageCreateChat 新建私聊消息
func MessageCreateChat(from int, to int, message string) {
	_, _ = DatabaseEngine.Table(new(MessageModel)).Insert(&MessageModel{
		Time:     time.Now(),
		Type:     "chat",
		FromUser: from,
		ToUser:   to,
		Message:  message,
	})
}

// MessageCreateToAdmins 新建对于所有admin用户的消息
func MessageCreateToAdmins(typeName string, message string) {
	var allAdminUsers []UserModel
	_ = DatabaseEngine.Table(new(UserModel)).Where("level >= ?", 2).Find(&allAdminUsers)
	for _, i := range allAdminUsers {
		_, _ = DatabaseEngine.Table(new(MessageModel)).Insert(&MessageModel{
			Time:     time.Now(),
			Type:     typeName,
			FromUser: 0,
			ToUser:   i.ID,
			Message:  message,
		})
	}
}
