package serviceImpl

import (
	"app/domain/service"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	oauth1Twitter "github.com/dghubble/oauth1/twitter"
	"os"
)

var loginConfig = oauth1.Config{
	ConsumerKey:    os.Getenv("twihaialert_app_consumer_key"),
	ConsumerSecret: os.Getenv("twihaialert_app_consumer_secret"),
	CallbackURL:    os.Getenv("twihaialert_app_host") + "/login/callback",
	Endpoint:       oauth1Twitter.AuthorizeEndpoint,
}

var addAccountConfig = oauth1.Config{
	ConsumerKey:    os.Getenv("twihaialert_app_consumer_key"),
	ConsumerSecret: os.Getenv("twihaialert_app_consumer_secret"),
	CallbackURL:    os.Getenv("twihaialert_app_host") + "/l/addAccount/callback",
	Endpoint:       oauth1Twitter.AuthorizeEndpoint,
}

type TwitterServiceImpl struct {
}

func (t TwitterServiceImpl) GetLoginAccountInfo(s string, s2 string) (int64, string, error) {
	token := oauth1.NewToken(s, s2)
	httpClient := loginConfig.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	u, _, err := client.Accounts.VerifyCredentials(nil)
	if err != nil {
		return 0, "", err
	}
	return u.ID, u.ScreenName, nil
}

func (t TwitterServiceImpl) GetLoginAccessToken(rt string, rs string, v string) (string, string, error) {
	return loginConfig.AccessToken(rt, rs, v)
}

func (t TwitterServiceImpl) GetLoginRequestConfig() (string, string, string, error) {
	rt, rs, err := loginConfig.RequestToken()
	if err != nil {
		return "", "", "", err
	}
	url, err := loginConfig.AuthorizationURL(rt)
	if err != nil {
		return "", "", "", err
	}
	return rt, rs, url.String(), nil
}

func (t TwitterServiceImpl) GetAddAccountInfo(s string, s2 string) (int64, string, error) {
	token := oauth1.NewToken(s, s2)
	httpClient := addAccountConfig.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	u, _, err := client.Accounts.VerifyCredentials(nil)
	if err != nil {
		return 0, "", err
	}
	return u.ID, u.ScreenName, nil
}

func (t TwitterServiceImpl) GetAddAccessToken(rt string, rs string, v string) (string, string, error) {
	return addAccountConfig.AccessToken(rt, rs, v)
}

func (t TwitterServiceImpl) GetAddRequestConfig() (string, string, string, error) {
	rt, rs, err := addAccountConfig.RequestToken()
	if err != nil {
		return "", "", "", err
	}
	url, err := addAccountConfig.AuthorizationURL(rt)
	if err != nil {
		return "", "", "", err
	}
	return rt, rs, url.String(), nil
}

func NewTwitterServiceImpl() service.TwitterService {
	return TwitterServiceImpl{}
}
