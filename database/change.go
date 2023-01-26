package database

import (
	"change/config"
	"time"
)

// ChangeCreateNew 创建一个新的change
func ChangeCreateNew(title string, description string, subject int, user int, want string) {
	_, _ = DatabaseEngine.Table(new(ChangeModel)).Insert(&ChangeModel{
		Time:        time.Now(),
		User:        user,
		Geter:       0,
		Subject:     subject,
		Title:       title,
		Description: description,
		Want:        want,
		State:       1,
	})
}

// ChangeCheckTime 检查用户是否能立刻创建一个change
func ChangeCheckTime(user int) bool {
	if UserGetLevelByID(user) > 1 {
		return true
	}
	if !UserHaveUserByID(user) {
		return false
	}
	var i ChangeModel
	_, _ = DatabaseEngine.Table(new(ChangeModel)).Where("user = ?", user).Desc("time").Get(&i)
	if time.Since(i.Time).Minutes() >= float64(config.ConfigData.Limit.Change) {
		return true
	}
	return false
}

// ChangeHaveByID 是否存在change
func ChangeHaveByID(id int) bool {
	cnt, _ := DatabaseEngine.Table(new(ChangeModel)).Where("id = ?", id).Count()
	return cnt != 0
}

// ChangeCheckUser 是否为change发布者
func ChangeCheckUser(id int, user int) bool {
	if !ChangeHaveByID(id) {
		return false
	}
	var c ChangeModel
	_, _ = DatabaseEngine.Table(new(ChangeModel)).Where("id = ?", id).Get(&c)
	return user == c.User
}
