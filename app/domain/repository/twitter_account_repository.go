package repository

import "app/domain/entity"

type TwitterAccountRepository interface {
	Create(account *entity.TwitterAccount) (*entity.TwitterAccount, error)
	Update(account *entity.TwitterAccount) (*entity.TwitterAccount, error)
	FindByTwitterId(twitterId int64) ([]entity.TwitterAccount, error)
	FindByUserId(userId int) ([]entity.TwitterAccount, error)
}
