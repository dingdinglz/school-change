package database

import "time"

type UserModel struct {
	ID       int       `xorm:"pk autoincr"` // id，唯一标识
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

type ApplyModel struct {
	Time     time.Time
	Username string
	Password string
	Ip       string
	Grade    int
	Class    int
	Realname string
}

func (i *ApplyModel) TableName() string {
	return "apply"
}

type SubjectModel struct {
	ID          int `xorm:"pk autoincr"`
	Name        string
	Description string
}

func (i *SubjectModel) TableName() string {
	return "subject"
}

type ChangeModel struct {
	ID          int `xorm:"pk autoincr"`
	Time        time.Time
	User        int
	Geter       int
	Subject     int
	Title       string
	Description string
	Want        string
	State       int // 1 未开始  2 交换中  3 交换结束
	Money       int // 价格
}

func (i *ChangeModel) TableName() string {
	return "change"
}

type MessageModel struct {
	Time     time.Time
	Type     string // chat为私聊
	FromUser int    //为0则是系统消息
	ToUser   int
	Message  string
}

func (i *MessageModel) TableName() string {
	return "message"
}

type ReportModel struct {
	Time    time.Time
	Change  int
	Message string
	User    int
}

func (i *ReportModel) TableName() string {
	return "report"
}
