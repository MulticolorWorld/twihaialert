package serviceImpl

import (
	"app/domain/service"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	oauth1Twitter "github.com/dghubble/oauth1/twitter"
	"os"
)

var config = oauth1.Config{
	ConsumerKey:    os.Getenv("twihaialert_app_consumer_key"),
	ConsumerSecret: os.Getenv("twihaialert_app_consumer_secret"),
	CallbackURL:    "http://localhost:1323/login/callback",
	Endpoint:       oauth1Twitter.AuthorizeEndpoint,
}

type TwitterServiceImpl struct {
}

func (t TwitterServiceImpl) GetAccountInfo(s string, s2 string) (int64, string, error) {
	token := oauth1.NewToken(s, s2)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	u, _, err := client.Accounts.VerifyCredentials(nil)
	if err != nil {
		return 0, "", err
	}
	return u.ID, u.ScreenName, nil
}

func (t TwitterServiceImpl) GetAccessToken(rt string, rs string, v string) (string, string, error) {
	return config.AccessToken(rt, rs, v)
}

func (t TwitterServiceImpl) GetRequestConfig() (string, string, string, error) {
	rt, rs, err := config.RequestToken()
	if err != nil {
		return "", "", "", err
	}
	url, err := config.AuthorizationURL(rt)
	if err != nil {
		return "", "", "", err
	}
	return rt, rs, url.String(), nil
}

func NewTwitterServiceImpl() service.TwitterService {
	return TwitterServiceImpl{}
}
