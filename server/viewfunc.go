package server

import (
	"change/database"
	"github.com/gofiber/template/html"
	"html/template"
)

func viewEngineAddFuncs(e *html.Engine) {
	e.AddFunc("User_GetRealNameByUsername", func(username string) template.HTML {
		return template.HTML(database.UserGetRealnameByUsername(username))
	})
	e.AddFunc("Subject_GetNameByID", func(id int) template.HTML {
		return template.HTML(database.SubjectGetNameByID(id))
	})
	e.AddFunc("User_GetUsernameByID", func(id int) template.HTML {
		var i database.UserModel
		_, _ = database.DatabaseEngine.Table(new(database.UserModel)).Where("id = ?", id).Get(&i)
		return template.HTML(i.Username)
	})
	e.AddFunc("User_GetRealnameByID", func(id int) template.HTML {
		var i database.UserModel
		_, _ = database.DatabaseEngine.Table(new(database.UserModel)).Where("id = ?", id).Get(&i)
		return template.HTML(i.Realname)
	})
}
