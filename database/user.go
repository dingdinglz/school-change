package database

import (
	"change/logger"
	"change/tools"
	"strconv"
	"time"
)

// UserGetNum 返回用户总数
func UserGetNum() int {
	res := GetSetting("user")
	res_i, _ := strconv.Atoi(res)
	return res_i
}

// UserGetNumReal 返回用户记录条数
func UserGetNumReal() int {
	res, _ := DatabaseEngine.Table(new(UserModel)).Count()
	return int(res)
}

// UserCreateNew 创建一个新用户，无任何过滤，纯表操作
func UserCreateNew(username string, password string, level int, ip string, grade int, class int, realname string) error {
	logger.ConsoleLogger.Debugln("创建新用户：" + username)
	_, err := DatabaseEngine.Table(new(UserModel)).Insert(UserModel{
		ID:       UserGetNum() + 1,
		Username: username,
		Password: tools.MD5(password),
		Level:    level,
		Ip:       ip,
		Time:     time.Now(),
		Grade:    grade,
		Class:    class,
		Realname: realname,
	})
	WriteOrUpdateSetting("user", strconv.Itoa(UserGetNum()+1))
	return err
}

// UserHaveUserByUsername 通过username判断用户是否存在
func UserHaveUserByUsername(username string) bool {
	cnt, _ := DatabaseEngine.Table(new(UserModel)).Where("username = ?", username).Count()
	return cnt != 0
}

// UserGetUserByUsername 通过username取用户
func UserGetUserByUsername(username string) UserModel {
	if !UserHaveUserByUsername(username) {
		return UserModel{}
	}
	var _User UserModel
	_, _ = DatabaseEngine.Table(new(UserModel)).Where("username = ?", username).Get(&_User)
	return _User
}

// UserCheck 检查username与password以及level/id是否配对
func UserCheck(username string, password string, level int, id int) bool {
	if !UserHaveUserByUsername(username) {
		return false
	}
	u := UserGetUserByUsername(username)
	if u.Password == password && u.Level == level && u.ID == id {
		return true
	}
	return false
}

// UserGetRealnameByUsername 根据username取realname
func UserGetRealnameByUsername(username string) string {
	if !UserHaveUserByUsername(username) {
		return ""
	}
	u := UserGetUserByUsername(username)
	return u.Realname
}

// UserHaveUserByID 根据id判断用户是否存在
func UserHaveUserByID(id int) bool {
	cnt, _ := DatabaseEngine.Table(new(UserModel)).Where("id = ?", id).Count()
	return cnt != 0
}

func UserGetLevelByID(id int) int {
	if !UserHaveUserByID(id) {
		return 0
	}
	var u UserModel
	_, _ = DatabaseEngine.Table(new(UserModel)).Where("id = ?", id).Get(&u)
	return u.Level
}

// UserCheckChangeAble 判断是否可以操作目标
func UserCheckChangeAble(myID int, opID int) bool {
	if myID == opID {
		return false
	}
	myLevel := UserGetLevelByID(myID)
	if myLevel == 3 {
		return true
	}
	opLevel := UserGetLevelByID(opID)
	if myLevel > opLevel {
		return true
	}
	return false
}
