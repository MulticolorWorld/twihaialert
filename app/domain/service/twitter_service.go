package service

type TwitterService interface {
	GetLoginRequestConfig() (string, string, string, error)
	GetLoginAccessToken(string, string, string) (string, string, error)
	GetLoginAccountInfo(string, string) (int64, string, error)
	GetAddRequestConfig() (string, string, string, error)
	GetAddAccessToken(string, string, string) (string, string, error)
	GetAddAccountInfo(string, string) (int64, string, error)
}
