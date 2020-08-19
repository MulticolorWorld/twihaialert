package service

import "app/domain/entity"

type TwitterService interface {
	GetRequestConfig(mode string) (rToken string, rSecret string, url string, err error)
	GetAccessToken(mode string, rToken string, rSecret string, v string) (aToken string, aSecret string, err error)
	GetAccountInfo(mode string, aToken string, aSecret string) (id int64, name string, err error)
	GetCountFromLastId(lastId int64, twitterId int64, aToken string, aSecret string) (count int, rtCount int, newLastId int64, err error)
	PostStatus(message string, dm int, accounts []entity.TwitterAccount) (err error)
}
