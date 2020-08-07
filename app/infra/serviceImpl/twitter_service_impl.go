package serviceImpl

import (
	"app/domain/service"
	"errors"
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

func getConfig(mode string) (c *oauth1.Config, err error) {
	if mode == "login" {
		return &loginConfig, nil
	}
	if mode == "addAccount" {
		return &addAccountConfig, nil
	}
	return nil, errors.New("mode invalid")
}

type TwitterServiceImpl struct {
}

func (t TwitterServiceImpl) GetRequestConfig(mode string) (rToken string, rSecret string, urlString string, err error) {
	config, _ := getConfig(mode)
	rToken, rSecret, err = config.RequestToken()
	if err != nil {
		return "", "", "", err
	}
	url, err := config.AuthorizationURL(rToken)
	if err != nil {
		return "", "", "", err
	}
	return rToken, rSecret, url.String(), nil
}

func (t TwitterServiceImpl) GetAccessToken(mode string, rToken string, rSecret string, v string) (aToken string, aSecret string, err error) {
	config, _ := getConfig(mode)
	return config.AccessToken(rToken, rSecret, v)
}

func (t TwitterServiceImpl) GetAccountInfo(mode string, aToken string, aSecret string) (id int64, name string, err error) {
	config, _ := getConfig(mode)
	token := oauth1.NewToken(aToken, aSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	u, _, err := client.Accounts.VerifyCredentials(nil)
	if err != nil {
		return 0, "", err
	}
	return u.ID, u.ScreenName, nil
}

func NewTwitterServiceImpl() service.TwitterService {
	return TwitterServiceImpl{}
}
