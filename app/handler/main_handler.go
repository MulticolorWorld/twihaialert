package handler

import (
	"app/useCase"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	oauth1Twitter "github.com/dghubble/oauth1/twitter"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

var config = oauth1.Config{
	ConsumerKey:    os.Getenv("twihaialert_app_consumer_key"),
	ConsumerSecret: os.Getenv("twihaialert_app_consumer_secret"),
	CallbackURL:    "http://localhost:1323/login/callback",
	Endpoint:       oauth1Twitter.AuthorizeEndpoint,
}

type MainHandler struct {
	mu useCase.MainUseCase
}

func NewMainHandler(mu useCase.MainUseCase) *MainHandler {
	return &MainHandler{mu: mu}
}

func (wh MainHandler) Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}

func (wh MainHandler) Login(c echo.Context) error {

	requestToken, requestSecret, _ := config.RequestToken()

	{
		sess, _ := session.Get("session", c)
		sess.Values["requestToken"] = requestToken
		sess.Values["requestSecret"] = requestSecret
		sess.Save(c.Request(), c.Response())
	}

	authorizationUrl, _ := config.AuthorizationURL(requestToken)
	return c.Redirect(http.StatusFound, authorizationUrl.String())
}

func (wh MainHandler) LoginCallback(c echo.Context) error {
	sess, _ := session.Get("session", c)
	requestSecret := sess.Values["requestSecret"].(string)
	requestToken, verifier, _ := oauth1.ParseAuthorizationCallback(c.Request())
	accessToken, accessSecret, _ := config.AccessToken(requestToken, requestSecret, verifier)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	user, _, _ := client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})
	return c.String(http.StatusOK, user.ScreenName+":"+user.Name)
}
