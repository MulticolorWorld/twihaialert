package useCase

import (
	"app/domain/entity"
	"app/domain/repository"
	"app/domain/service"
	"app/errors"
	"strconv"
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
	return mu.ts.GetLoginRequestConfig()
}

func (mu MainUseCase) Login(rt string, rs string, v string) (int, error) {
	at, as, err := mu.ts.GetLoginAccessToken(rt, rs, v)
	if err != nil {
		return 0, err
	}
	id, name, err := mu.ts.GetLoginAccountInfo(at, as)
	if err != nil {
		return 0, err
	}
	tas, err := mu.tar.FindByTwitterId(id)
	if err != nil {
		return 0, err
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
			return 0, err
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

func (mu MainUseCase) UpdateConfig(d string, id int) error {
	us, err := mu.ur.FindById(id)
	if err != nil {
		return err
	}
	u := &us[0]
	dn, err := strconv.Atoi(d)
	if err != nil {
		return err
	}
	u.DMNotification = dn
	_, err = mu.ur.Update(u)
	if err != nil {
		return err
	}
	return nil
}

func (mu MainUseCase) PreAddAccount() (string, string, string, error) {
	return mu.ts.GetAddRequestConfig()
}

func (mu MainUseCase) AddAccount(rt string, rs string, v string, userId int) error {
	at, as, err := mu.ts.GetAddAccessToken(rt, rs, v)
	if err != nil {
		return err
	}
	id, name, err := mu.ts.GetAddAccountInfo(at, as)
	if err != nil {
		return err
	}
	tas, err := mu.tar.FindByTwitterId(id)
	if err != nil {
		return err
	}
	if len(tas) != 0 {
		return &errors.AccountAlreadyExistError{}
	} else {
		ta := entity.NewTwitterAccount()
		ta.TwitterId = id
		ta.UserId = userId
		ta.ScreenName = name
		ta.AccessToken = at
		ta.AccessTokenSecret = as
		ta, err = mu.tar.Create(ta)
		if err != nil {
			return err
		}
	}
	return nil
}
