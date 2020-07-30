package serviceImpl

import "app/domain/service"

type TwitterServiceImpl struct {
}

func NewTwitterServiceImpl() service.TwitterService {
	return TwitterServiceImpl{}
}
