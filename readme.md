2022年安徽省信息素养大赛一等奖作品

整理仓库时感慨时间之快，放出当时青涩的代码供大家参考。

# 校园旧物交换平台

![](./web/public/logo.png)

## 作者：合肥市第七中学dinglz

## 介绍

为防止浪费，使各类物品充分利用，方便学校学生的物品交换。开发了这个校园旧物交换平台。灵感来自于经常能看到高三学长学姐出售书籍等物品。

## 功能

本项目实现了一个在线网站，用于学生交换、出售旧物。通过本平台，学生可以上传需要交换或出售的物品，相互交流。

## 架构

本项目分别设计了前台和后台。

### 前台

前台允许学生使用最基本的功能。

### 后台

后台允许校方管理员管理平台，处理用户申请，处理举报，管理用户，接受反馈等。

## 技术选型

| 名称      | 选型                                                                |
|:-------:|:-----------------------------------------------------------------:|
| 开发语言    | golang                                                            |
| 数据库     | sqlite                                                            |
| web后端框架 | [fiber](https://github.com/gofiber/fiber)                         |
| 前端框架    | [bulma](https://bulma.io/) && [layui](https://layui.gitee.io/v2/) |
| orm框架   | [xorm](https://xorm.io/zh/)                                       |

## 编译

```shell
go build
```

得到的可执行文件与web目录（包含网页文件，css，js以及模板文件）即为最终的可用程序。

## 运行

初次运行时，会沿用默认参数，并在目录下生成data目录。其中的change.config文件需要在使用者修改后并重新启动。change.config的参数如下。具体运行流程和功能介绍参见视频。

### change.config文件格式

文件为json格式，对应的功能如表格。

| 名称     | 对应功能                                |
| ------ | ----------------------------------- |
| port   | 启动的端口。正常http对应域名端口为80，ssl下https为443 |
| ssl    | 是否启动ssl，如果启动，则以https形式支持            |
| debug  | 是否启动调试模式。                           |
| school | 使用的学校，例如，在合肥市第七中学内使用则为合肥市第七中学       |
| limit  | 限制块，详情见下                            |

#### limit格式

| 名称     | 对应功能                     |
| ------ | ------------------------ |
| change | 学生每次发布交换必须间隔的时间，防止恶意批量发布 |
| web    | 每秒钟个人的访问次数的限制，防止恶意访问     |
