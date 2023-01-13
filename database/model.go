package database

import "time"

type UserModel struct {
	ID       int       // id，唯一标识
	Username string    // 用户名
	Password string    // 密码，md5形式
	Level    int       // 权限等级 1学生，2老师，3后台管理
	Ip       string    // Ip
	Time     time.Time // 注册时间
	Grade    int       // 年级
	Class    int       // 班级
	Realname string    // 真实姓名
}

func (i *UserModel) TableName() string {
	return "user"
}

type SettingModel struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (i *SettingModel) TableName() string {
	return "setting"
}
