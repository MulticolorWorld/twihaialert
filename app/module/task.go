package module

import (
	"app/domain/entity"
	"app/infra/persistence"
	"app/infra/serviceImpl"
	"github.com/labstack/gommon/log"
	"github.com/panjf2000/ants/v2"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var service = serviceImpl.NewTwitterServiceImpl()

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
	now := time.Now()

	for i, a := range accounts {
		c, rc, newId, err := service.GetCountFromLastId(a.LastTweetId, a.TwitterId, a.AccessToken, a.AccessTokenSecret)
		if err != nil {
			log.Error("Account:" + strconv.Itoa(a.ID) + " failed")
			continue
		}
		accounts[i].DailyCount += c
		accounts[i].DailyCountRt += rc
		accounts[i].LastTweetId = newId
		accounts[i].CountUpdate = &now
	}

	if user.LastNotify == nil {
		user.LastNotify = &now
	} else if user.LastNotify.YearDay() != now.YearDay() {
		count := 0
		rtCount := 0
		head := ""
		for _, a := range accounts {
			count += a.DailyCount
			rtCount += a.DailyCountRt
			head += "@" + a.ScreenName + " "
		}
		message := head + user.LastNotify.Format("01-02") + "のポスト数：" + strconv.Itoa(count) + " (うちRT：" + strconv.Itoa(rtCount) + ")"
		err := service.PostStatus(message, user.DMNotification, accounts)
		if err != nil {
			return nil, nil, err
		}
		user.LastNotifyCount = 0
		user.LastNotify = &now
		for i := range accounts {
			accounts[i].DailyCount = 0
			accounts[i].DailyCountRt = 0
		}
	} else {
		interval := 100
		count := 0
		head := ""
		for _, a := range accounts {
			count += a.DailyCount
			head += "@" + a.ScreenName + " "
		}
		if count/interval > user.LastNotifyCount/interval {
			message := head + user.LastNotify.Format("01-02") + "に入ってから" + strconv.Itoa((count/interval)*interval) +
				"ポストに到達しました！ (現在およそ" + strconv.Itoa(count) + ")"
			err := service.PostStatus(message, user.DMNotification, accounts)
			if err != nil {
				return nil, nil, err
			}
			user.LastNotifyCount = (count / interval) * interval
			user.LastNotify = &now
		}
	}
	return user, accounts, nil
}
