package database

import (
	"change/logger"
	"change/tools"
	"strconv"
	"time"
)

func UserGetNum() int {
	res := GetSetting("user")
	res_i, _ := strconv.Atoi(res)
	return res_i
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
