package handler

import (
	"app/domain/entity"
	"app/useCase"
	"github.com/dghubble/oauth1"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

type MainHandler struct {
	mu useCase.MainUseCase
}

type Data struct {
	User     *entity.User
	Accounts []entity.TwitterAccount
}

func NewMainHandler(mu useCase.MainUseCase) *MainHandler {
	return &MainHandler{mu: mu}
}

func (mh MainHandler) Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func (mh MainHandler) Login(c echo.Context) error {
	rt, rs, url, _ := mh.mu.PreLogin()

	sess, _ := session.Get("session", c)
	sess.Values["requestToken"] = rt
	sess.Values["requestSecret"] = rs
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, url)
}

func (mh MainHandler) LoginCallback(c echo.Context) error {
	sess, _ := session.Get("session", c)
	rs := sess.Values["requestSecret"].(string)
	rt, v, err := oauth1.ParseAuthorizationCallback(c.Request())
	if err != nil {
		return err
	}
	id, err := mh.mu.Login(rt, rs, v)
	if err != nil {
		return err
	}
	sess.Values["userId"] = id
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, "/myPage")
}

func (mh MainHandler) MyPage(c echo.Context) error {
	sess, _ := session.Get("session", c)
	userId := sess.Values["userId"].(int)
	u, tas, err := mh.mu.FindUserInfo(userId)
	if err != nil {
		return err
	}
	data := Data{
		User:     u,
		Accounts: tas,
	}
	return c.Render(http.StatusOK, "myPage", data)
}
