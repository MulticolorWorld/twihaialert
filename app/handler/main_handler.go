package handler

import (
	"app/domain/entity"
	errors2 "app/errors"
	"app/useCase"
	"errors"
	"github.com/dghubble/oauth1"
	"github.com/google/uuid"
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
	Token    string
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

	return c.Redirect(http.StatusFound, "/l/myPage")
}

func (mh MainHandler) MyPage(c echo.Context) error {
	sess, _ := session.Get("session", c)
	userId := sess.Values["userId"].(int)
	u, tas, err := mh.mu.FindUserInfo(userId)
	if err != nil {
		return err
	}
	d := Data{
		User:     u,
		Accounts: tas,
		Token:    "",
	}
	return c.Render(http.StatusOK, "myPage", d)
}

func (mh MainHandler) ConfigInput(c echo.Context) error {
	sess, _ := session.Get("session", c)
	userId := sess.Values["userId"].(int)
	u, tas, err := mh.mu.FindUserInfo(userId)
	if err != nil {
		return err
	}
	t := uuid.New().String()
	d := Data{
		User:     u,
		Accounts: tas,
		Token:    t,
	}
	sess.Values["csrfToken"] = t
	sess.Save(c.Request(), c.Response())
	return c.Render(http.StatusOK, "configInput", d)
}

func (mh MainHandler) Config(c echo.Context) error {
	sess, _ := session.Get("session", c)
	st := sess.Values["csrfToken"].(string)
	ft := c.FormValue("token")
	if st != ft {
		return c.Redirect(http.StatusFound, "/error/wrongToken")
	}
	dn := c.FormValue("dmNotification")
	userId := sess.Values["userId"].(int)
	err := mh.mu.UpdateConfig(dn, userId)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, "/l/configFinish")
}

func (mh MainHandler) ConfigFinish(c echo.Context) error {
	return c.Render(http.StatusOK, "configFinish", nil)
}

func (mh MainHandler) AddAccount(c echo.Context) error {
	rt, rs, url, _ := mh.mu.PreAddAccount()

	sess, _ := session.Get("session", c)
	sess.Values["requestToken"] = rt
	sess.Values["requestSecret"] = rs
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, url)
}

func (mh MainHandler) AddAccountCallback(c echo.Context) error {
	sess, _ := session.Get("session", c)
	rs := sess.Values["requestSecret"].(string)
	rt, v, err := oauth1.ParseAuthorizationCallback(c.Request())
	if err != nil {
		return err
	}
	userId := sess.Values["userId"].(int)
	err = mh.mu.AddAccount(rt, rs, v, userId)
	if err != nil {
		if errors.Is(err, &errors2.AccountAlreadyExistError{}) {
			return c.Redirect(http.StatusFound, "/error/accountAlreadyExist")
		}
		return err
	}
	return c.Redirect(http.StatusFound, "/l/myPage")
}
