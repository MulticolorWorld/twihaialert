package service

type TwitterService interface {
	GetRequestConfig() (string, string, string, error)
	GetAccessToken(string, string, string) (string, string, error)
	GetAccountInfo(string, string) (int64, string, error)
}
