package main

import (
	"app/handler"
	"app/infra/persistence"
	"app/infra/serviceImpl"
	"app/useCase"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"html/template"
	"io"
)

func main() {
	db, err := persistence.InitDBConnection()
	if err != nil {
		panic("DB接続エラー")
	}
	defer db.Close()

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.gohtml")),
	}

	up := persistence.NewUserPersistence(db)
	tap := persistence.NewTwitterAccountPersistence(db)
	ts := serviceImpl.NewTwitterServiceImpl()
	mu := useCase.NewMainUseCase(up, tap, ts)
	mh := handler.NewMainHandler(*mu)

	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore(securecookie.GenerateRandomKey(32))))
	e.Renderer = t
	e.Logger.SetLevel(log.DEBUG)

	e.GET("/", mh.Index)
	e.POST("/login", mh.Login)
	e.GET("/login/callback", mh.LoginCallback)
	e.GET("/myPage", mh.MyPage)

	e.Logger.Fatal(e.Start(":1323"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
