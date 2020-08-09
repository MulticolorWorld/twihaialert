package main

import (
	"app/domain/entity"
	"app/infra/persistence"
	"fmt"
	"github.com/pkg/profile"
)

func Task() {
	defer profile.Start().Stop()

	db, err := persistence.InitDBConnection()
	if err != nil {
		panic("DB接続エラー")
	}
	defer db.Close()

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
	fmt.Println(len(users))
	fmt.Println(len(accountMap))
}
