package service

type TwitterService interface {
	GetRequestConfig(mode string) (rToken string, rSecret string, url string, err error)
	GetAccessToken(mode string, rToken string, rSecret string, v string) (aToken string, aSecret string, err error)
	GetAccountInfo(mode string, aToken string, aSecret string) (id int64, name string, err error)
}
