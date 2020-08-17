package module

import (
	"app/domain/entity"
	"app/infra/persistence"
	"errors"
	"github.com/panjf2000/ants/v2"
	"net/http"
	"sync"
)

func Task() {
	db, err := persistence.InitDBConnection()
	if err != nil {
		panic("DB接続エラー")
	}
	db.DB().SetMaxIdleConns(110)
	db.DB().SetMaxOpenConns(110)
	db.LogMode(true)
	defer db.Close()

	http.DefaultTransport.(*http.Transport).MaxIdleConns = 0
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 0

	userPer := persistence.NewUserPersistence(db)
	accountPer := persistence.NewTwitterAccountPersistence(db)

	users, err := userPer.FindAll()
	if err != nil {
		panic("user全件取得時エラー")
	}
	accounts, err := accountPer.FindAll()
	if err != nil {
		panic("account全件取得時エラー")
	}
	accountMap := make(map[int][]entity.TwitterAccount)
	for _, a := range accounts {
		s, ok := accountMap[a.UserId]
		if ok {
			s = append(s, a)
		} else {
			s = []entity.TwitterAccount{a}
		}
		accountMap[a.UserId] = s
	}

	p, _ := ants.NewPool(500)
	defer p.Release()
	var wg sync.WaitGroup
	for _, val := range users {
		wg.Add(1)
		user := val
		_ = p.Submit(func() {
			defer wg.Done()
			updatedUser, updatedAccounts, err := oneUserTask(&user, accountMap[user.ID])
			if err != nil {
				return
			}
			_, err = userPer.Update(updatedUser)
			if err != nil {
				return
			}
			for _, a := range updatedAccounts {
				_, err = accountPer.Update(&a)
			}
			if err != nil {
				return
			}
			return
		})
	}
	wg.Wait()
}

func oneUserTask(user *entity.User, accounts []entity.TwitterAccount) (*entity.User, []entity.TwitterAccount, error) {
	resp, err := http.Get("https://www.google.com/")
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	return user, accounts, errors.New("test")
}
