package persistence

import (
	"app/domain/entity"
	"app/domain/repository"
	"github.com/jinzhu/gorm"
)

type TwitterAccountPersistence struct {
	db *gorm.DB
}

func (t TwitterAccountPersistence) FindAll() ([]entity.TwitterAccount, error) {
	var accounts []entity.TwitterAccount
	err := t.db.Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (t TwitterAccountPersistence) Delete(account *entity.TwitterAccount) error {
	err := t.db.Delete(account).Error
	if err != nil {
		return err
	}
	return nil
}

func (t TwitterAccountPersistence) Create(account *entity.TwitterAccount) (*entity.TwitterAccount, error) {
	err := t.db.Create(account).Error
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (t TwitterAccountPersistence) Update(account *entity.TwitterAccount) (*entity.TwitterAccount, error) {
	err := t.db.Save(account).Error
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (t TwitterAccountPersistence) FindByTwitterId(twitterId int64) ([]entity.TwitterAccount, error) {
	var accounts []entity.TwitterAccount
	err := t.db.Where("twitter_id = ?", twitterId).Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (t TwitterAccountPersistence) FindByUserId(userId int) ([]entity.TwitterAccount, error) {
	var accounts []entity.TwitterAccount
	err := t.db.Where("user_id = ?", userId).Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func NewTwitterAccountPersistence(db *gorm.DB) repository.TwitterAccountRepository {
	return &TwitterAccountPersistence{db: db}
}
