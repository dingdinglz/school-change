package main

import (
	"change/config"
	"fmt"
)

func main() {
	fmt.Println("校园旧物交换平台")
	fmt.Println("作者：合肥市第七中学丁励治")
	fmt.Println("鸣谢：感谢老师，各位开源作者，以及帮助测试本项目的各位同学！没有你们就没有这个项目。")
	InitPath()
	config.InitConfig()
	InitLogger()
	ConsoleLogger.Infoln("初始化已完成！")
	FileLogger.Infoln("服务已启动！")
}
