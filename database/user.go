package database

import (
	"change/logger"
	"change/tools"
	"errors"
	"time"
)

// UserGetNumReal 返回用户记录条数
func UserGetNumReal() int {
	res, _ := DatabaseEngine.Table(new(UserModel)).Count()
	return int(res)
}

// UserCreateNew 创建一个新用户，无任何过滤，纯表操作
func UserCreateNew(username string, password string, level int, ip string, grade int, class int, realname string) error {
	logger.ConsoleLogger.Debugln("创建新用户：" + username)
	_, err := DatabaseEngine.Table(new(UserModel)).Insert(UserModel{
		Username: username,
		Password: tools.MD5(password),
		Level:    level,
		Ip:       ip,
		Time:     time.Now(),
		Grade:    grade,
		Class:    class,
		Realname: realname,
	})
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

// UserCheckCreateAble 判断用户是否可以被创建
func UserCheckCreateAble(username string, password string, level int, class int, grade int) (bool, error) {
	if len(username) < 6 || len(username) > 15 {
		return false, errors.New("用户名长度不符合规则！")
	}
	if !tools.IsAllEnglish(username) {
		return false, errors.New("用户名必须只包含英文字母！")
	}
	if UserHaveUserByUsername(username) {
		return false, errors.New("用户名已被使用！")
	}
	if len(password) < 8 || len(password) > 50 {
		return false, errors.New("密码长度不符合！")
	}
	if level < 1 || level > 3 {
		return false, errors.New("level不符合规则！")
	}
	if class <= 0 || grade <= 0 {
		return false, errors.New("班级或年级有误！")
	}
	return true, nil
}

// UserApplyCreate 创建用户注册申请
func UserApplyCreate(username string, password string, ip string, grade int, class int, realname string) {
	_, _ = DatabaseEngine.Table(new(ApplyModel)).Insert(&ApplyModel{
		Time:     time.Now(),
		Username: username,
		Password: password,
		Ip:       ip,
		Grade:    grade,
		Class:    class,
		Realname: realname,
	})
}

// UserApplyHaveIP 用户ip是否已经申请过
func UserApplyHaveIP(ip string) bool {
	cnt, _ := DatabaseEngine.Table(new(ApplyModel)).Where("ip = ?", ip).Count()
	return cnt != 0
}

// UserApplyPass 通过申请
func UserApplyPass(ip string) error {
	if !UserApplyHaveIP(ip) {
		return errors.New("申请不存在！")
	}
	var u ApplyModel
	_, _ = DatabaseEngine.Table(new(ApplyModel)).Where("ip = ?", ip).Get(&u)
	cCreate, err := UserCheckCreateAble(u.Username, u.Password, 1, u.Class, u.Grade)
	if !cCreate {
		return err
	}
	_ = UserCreateNew(u.Username, u.Password, 1, u.Ip, u.Grade, u.Class, u.Realname)
	_, _ = DatabaseEngine.Table(new(ApplyModel)).Where("ip = ?", ip).Delete()
	return nil
}

func UserApplyStop(ip string) {
	_, _ = DatabaseEngine.Table(new(ApplyModel)).Where("ip = ?", ip).Delete()
}

// UserGetRealnameByID 根据id取realname
func UserGetRealnameByID(id int) string {
	if !UserHaveUserByID(id) {
		return ""
	}
	var u UserModel
	_, _ = DatabaseEngine.Table(new(UserModel)).Where("id = ?", id).Get(&u)
	return u.Realname
}
