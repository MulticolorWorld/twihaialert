package useCase

import (
	"app/domain/entity"
	"app/domain/repository"
	"app/domain/service"
	"time"
)

type MainUseCase struct {
	ur  repository.UserRepository
	tar repository.TwitterAccountRepository
	ts  service.TwitterService
}

func NewMainUseCase(ur repository.UserRepository, tar repository.TwitterAccountRepository, ts service.TwitterService) *MainUseCase {
	return &MainUseCase{ur: ur, tar: tar, ts: ts}
}

func (mu MainUseCase) PreLogin() (string, string, string, error) {
	return mu.ts.GetRequestConfig()
}

func (mu MainUseCase) Login(rt string, rs string, v string) (int, error) {
	at, as, err := mu.ts.GetAccessToken(rt, rs, v)
	if err != nil {
		return 0, nil
	}
	id, name, err := mu.ts.GetAccountInfo(at, as)
	if err != nil {
		return 0, nil
	}
	tas, err := mu.tar.FindByTwitterId(id)
	if err != nil {
		return 0, nil
	}
	if len(tas) != 0 { //ログイン
		ta := &tas[0]
		ta.ScreenName = name
		ta, err = mu.tar.Update(ta)
		if err != nil {
			return 0, nil
		}
		us, err := mu.ur.FindById(ta.UserId)
		if err != nil {
			return 0, nil
		}
		u := &us[0]
		now := time.Now()
		u.LastLogin = &now
		u, err = mu.ur.Update(u)
		if err != nil {
			return 0, nil
		}
		return u.ID, nil
	} else { //新規登録
		u := entity.NewUser()
		u, err = mu.ur.Create(u)
		if err != nil {
			return 0, nil
		}
		ta := entity.NewTwitterAccount()
		ta.TwitterId = id
		ta.UserId = u.ID
		ta.ScreenName = name
		ta.AccessToken = at
		ta.AccessTokenSecret = as
		ta, err = mu.tar.Create(ta)
		if err != nil {
			return 0, nil
		}
		return u.ID, nil
	}
}

func (mu MainUseCase) FindUserInfo(userId int) (*entity.User, []entity.TwitterAccount, error) {
	us, err := mu.ur.FindById(userId)
	if err != nil {
		return nil, nil, err
	}
	u := &us[0]
	tas, err := mu.tar.FindByUserId(u.ID)
	if err != nil {
		return nil, nil, err
	}
	return u, tas, nil
}
