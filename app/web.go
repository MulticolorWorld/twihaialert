package main

import (
	"app/handler"
	"app/infra/persistence"
	"app/infra/serviceImpl"
	"app/useCase"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"html/template"
	"io"
	"net/http"
)

func Web() {
	db, err := persistence.InitDBConnection()
	if err != nil {
		panic("DB接続エラー")
	}
	db.LogMode(true)
	defer db.Close()

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.gohtml")),
	}

	up := persistence.NewUserPersistence(db)
	tap := persistence.NewTwitterAccountPersistence(db)
	ts := serviceImpl.NewTwitterServiceImpl()
	mu := useCase.NewMainUseCase(up, tap, ts)
	mh := handler.NewMainHandler(*mu)
	eh := handler.NewErrorHandler()

	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore(securecookie.GenerateRandomKey(32))))
	e.Renderer = t
	e.Static("/", "public/assets")
	e.Logger.SetLevel(log.DEBUG)

	e.GET("/", mh.Index)
	e.GET("/login", mh.Login)
	e.GET("/login/callback", mh.LoginCallback)

	l := e.Group("/l")
	l.Use(loginCheckMiddleware)
	l.GET("/myPage", mh.MyPage)
	l.GET("/configInput", mh.ConfigInput)
	l.POST("/config", mh.Config)
	l.GET("/configFinish", mh.ConfigFinish)

	er := e.Group("/error")
	er.GET("/wrongToken", eh.WrongToken)

	e.Logger.Fatal(e.Start(":1323"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func loginCheckMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		userId := sess.Values["userId"]
		if userId == nil {
			return c.Redirect(http.StatusFound, "/")
		}
		err := next(c)
		return err
	}
}
