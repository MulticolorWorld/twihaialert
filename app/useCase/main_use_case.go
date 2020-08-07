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
	userRepository    repository.UserRepository
	accountRepository repository.TwitterAccountRepository
	twitterService    service.TwitterService
}

func NewMainUseCase(ur repository.UserRepository, tar repository.TwitterAccountRepository, ts service.TwitterService) *MainUseCase {
	return &MainUseCase{userRepository: ur, accountRepository: tar, twitterService: ts}
}

func (u MainUseCase) PreLogin() (string, string, string, error) {
	return u.twitterService.GetRequestConfig("login")
}

func (u MainUseCase) Login(rToken string, rSecret string, v string) (id int, err error) {
	aToken, aSecret, err := u.twitterService.GetAccessToken("login", rToken, rSecret, v)
	if err != nil {
		return 0, err
	}
	twitterId, name, err := u.twitterService.GetAccountInfo("login", aToken, aSecret)
	if err != nil {
		return 0, err
	}
	tas, err := u.accountRepository.FindByTwitterId(twitterId)
	if err != nil {
		return 0, err
	}
	if len(tas) != 0 { //ログイン
		ta := &tas[0]
		ta.ScreenName = name
		ta, err = u.accountRepository.Update(ta)
		if err != nil {
			return 0, err
		}
		users, err := u.userRepository.FindById(ta.UserId)
		if err != nil {
			return 0, err
		}
		user := &users[0]
		now := time.Now()
		user.LastLogin = &now
		user, err = u.userRepository.Update(user)
		if err != nil {
			return 0, err
		}
		return user.ID, nil
	} else { //新規登録
		user := entity.NewUser()
		user, err = u.userRepository.Create(user)
		if err != nil {
			return 0, err
		}
		ta := entity.NewTwitterAccount()
		ta.TwitterId = twitterId
		ta.UserId = user.ID
		ta.ScreenName = name
		ta.AccessToken = aToken
		ta.AccessTokenSecret = aSecret
		ta, err = u.accountRepository.Create(ta)
		if err != nil {
			return 0, err
		}
		return user.ID, nil
	}
}

func (u MainUseCase) FindUserInfo(userId int) (*entity.User, []entity.TwitterAccount, error) {
	users, err := u.userRepository.FindById(userId)
	if err != nil {
		return nil, nil, err
	}
	user := &users[0]
	accounts, err := u.accountRepository.FindByUserId(user.ID)
	if err != nil {
		return nil, nil, err
	}
	return user, accounts, nil
}

func (u MainUseCase) UpdateConfig(s string, id int) error {
	users, err := u.userRepository.FindById(id)
	if err != nil {
		return err
	}
	user := &users[0]
	notify, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	user.DMNotification = notify
	_, err = u.userRepository.Update(user)
	if err != nil {
		return err
	}
	return nil
}

func (u MainUseCase) PreAddAccount() (string, string, string, error) {
	return u.twitterService.GetRequestConfig("addAccount")
}

func (u MainUseCase) AddAccount(rToken string, rSecret string, v string, userId int) error {
	aToken, aSecret, err := u.twitterService.GetAccessToken("addAccount", rToken, rSecret, v)
	if err != nil {
		return err
	}
	id, name, err := u.twitterService.GetAccountInfo("addAccount", aToken, aSecret)
	if err != nil {
		return err
	}
	accounts, err := u.accountRepository.FindByTwitterId(id)
	if err != nil {
		return err
	}
	if len(accounts) != 0 {
		return &errors.AccountAlreadyExistError{}
	} else {
		account := entity.NewTwitterAccount()
		account.TwitterId = id
		account.UserId = userId
		account.ScreenName = name
		account.AccessToken = aToken
		account.AccessTokenSecret = aSecret
		account, err = u.accountRepository.Create(account)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u MainUseCase) Remove(id int) error {
	users, err := u.userRepository.FindById(id)
	if err != nil {
		return err
	}
	user := users[0]
	accounts, err := u.accountRepository.FindByUserId(id)
	if err != nil {
		return err
	}
	for _, a := range accounts {
		err = u.accountRepository.Delete(&a)
		if err != nil {
			return err
		}
	}
	err = u.userRepository.Delete(&user)
	if err != nil {
		return err
	}
	return nil
}
