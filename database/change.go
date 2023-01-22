package database

import "time"

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
	var i ChangeModel
	_, _ = DatabaseEngine.Table(new(ChangeModel)).Where("user = ?", user).Desc("time").Get(&i)
	if time.Since(i.Time).Minutes() >= 30 {
		return true
	}
	return false
}
