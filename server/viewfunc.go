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
}
