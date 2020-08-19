package serviceImpl

import (
	"app/domain/entity"
	"app/domain/service"
	"errors"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	oauth1Twitter "github.com/dghubble/oauth1/twitter"
	"os"
	"strconv"
	"time"
)

var baseConfig = oauth1.Config{
	ConsumerKey:    os.Getenv("twihaialert_app_consumer_key"),
	ConsumerSecret: os.Getenv("twihaialert_app_consumer_secret"),
	CallbackURL:    "",
	Endpoint:       oauth1Twitter.AuthorizeEndpoint,
}

func getConfig(mode string) (c *oauth1.Config, err error) {
	config := baseConfig
	if mode == "login" {
		config.CallbackURL = os.Getenv("twihaialert_app_host") + "/login/callback"
		return &config, nil
	}
	if mode == "addAccount" {
		config.CallbackURL = os.Getenv("twihaialert_app_host") + "/l/addAccount/callback"
		return &config, nil
	}
	return nil, errors.New("mode invalid")
}

type TwitterServiceImpl struct {
}

func (t TwitterServiceImpl) PostStatus(message string, dm int, accounts []entity.TwitterAccount) (err error) {
	config := baseConfig
	for _, a := range accounts {
		token := oauth1.NewToken(a.AccessToken, a.AccessTokenSecret)
		httpClient := config.Client(oauth1.NoContext, token)
		client := twitter.NewClient(httpClient)
		if dm == 0 {
			_, _, err := client.Statuses.Update(message, nil)
			if err != nil {
				continue
			}
		} else {
			_, _, err := client.DirectMessages.EventsNew(&twitter.DirectMessageEventsNewParams{
				Event: &twitter.DirectMessageEvent{
					Type: "message_create",
					Message: &twitter.DirectMessageEventMessage{
						Target: &twitter.DirectMessageTarget{
							RecipientID: strconv.FormatInt(a.TwitterId, 10),
						},
						Data: &twitter.DirectMessageData{
							Text: message,
						},
					},
				},
			})
			if err != nil {
				continue
			}
		}
		return nil
	}
	return errors.New("User:" + strconv.Itoa(accounts[0].UserId) + " not notify.")
}

func (t TwitterServiceImpl) GetCountFromLastId(lastId int64, twitterId int64, aToken string, aSecret string) (count int, rtCount int, newLastId int64, err error) {
	config := baseConfig
	token := oauth1.NewToken(aToken, aSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	now := time.Now()
	bod := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tweetList := make([]twitter.Tweet, 0)

	if lastId == 0 {
		tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
			UserID: twitterId,
		})
		if err != nil {
			return 0, 0, 0, err
		}
		var c = 0
		for _, t := range tweets {
			st, _ := t.CreatedAtTime()
			if st.After(bod) {
				tweetList = append(tweetList, t)
				c++
			}
		}
		if c == 0 {
			return 0, 0, lastId, nil
		}
		maxId := tweetList[len(tweetList)-1].ID - 1
		for {
			tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
				UserID: twitterId,
				MaxID:  maxId,
			})
			if err != nil {
				return 0, 0, 0, err
			}
			var c = 0
			for _, t := range tweets {
				st, _ := t.CreatedAtTime()
				if st.After(bod) {
					tweetList = append(tweetList, t)
					c++
				}
			}
			if c == 0 {
				break
			}
			maxId = tweetList[len(tweetList)-1].ID - 1
		}
	} else {
		tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
			UserID:  twitterId,
			SinceID: lastId,
		})
		if err != nil {
			return 0, 0, 0, err
		}
		var c = 0
		for _, t := range tweets {
			st, _ := t.CreatedAtTime()
			if st.After(bod) {
				tweetList = append(tweetList, t)
				c++
			}
		}
		if c == 0 {
			return 0, 0, lastId, nil
		}
		maxId := tweetList[len(tweetList)-1].ID - 1
		for {
			tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
				UserID:  twitterId,
				MaxID:   maxId,
				SinceID: lastId,
			})
			if err != nil {
				return 0, 0, 0, err
			}
			var c = 0
			for _, t := range tweets {
				st, _ := t.CreatedAtTime()
				if st.After(bod) {
					tweetList = append(tweetList, t)
					c++
				}
			}
			if c == 0 {
				break
			}
			maxId = tweetList[len(tweetList)-1].ID - 1
		}
	}

	count = 0
	rtCount = 0
	for _, t := range tweetList {
		count += 1
		if t.Retweeted {
			rtCount += 1
		}
	}
	newLastId = tweetList[0].ID
	return count, rtCount, newLastId, nil
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
