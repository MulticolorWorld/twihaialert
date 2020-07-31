package entity

import "time"

type TwitterAccount struct {
	ID                int
	UserId            int
	TwitterId         int64
	ScreenName        string
	AccessToken       string
	AccessTokenSecret string
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
	LastTweetId       int64
	DailyCount        int
	DailyCountRt      int
	CountUpdate       *time.Time
}

func NewTwitterAccount() *TwitterAccount {
	now := time.Now()

	ta := new(TwitterAccount)
	ta.ID = 0
	ta.UserId = 0
	ta.TwitterId = 0
	ta.ScreenName = ""
	ta.AccessToken = ""
	ta.AccessTokenSecret = ""
	ta.CreatedAt = &now
	ta.UpdatedAt = &now
	ta.LastTweetId = 0
	ta.DailyCount = 0
	ta.DailyCountRt = 0
	ta.CountUpdate = nil
	return ta
}

func (ta TwitterAccount) TableName() string {
	return "twitter_account"
}
