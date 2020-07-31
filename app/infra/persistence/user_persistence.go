package persistence

import (
	"app/domain/entity"
	"app/domain/repository"
	"github.com/jinzhu/gorm"
	"os"
)

type UserPersistence struct {
	db *gorm.DB
}

func (up UserPersistence) Create(u *entity.User) (*entity.User, error) {
	err := up.db.Create(u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (up UserPersistence) Update(u *entity.User) (*entity.User, error) {
	err := up.db.Save(u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (up UserPersistence) FindById(id int) ([]entity.User, error) {
	var users []entity.User
	err := up.db.Where("id = ?", id).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func NewUserPersistence(db *gorm.DB) repository.UserRepository {
	return &UserPersistence{db: db}
}

func InitDBConnection() (*gorm.DB, error) {
	user := os.Getenv("twihaialert_db_user")
	pass := os.Getenv("twihaialert_db_password")
	name := os.Getenv("twihaialert_db_name")

	db, err := gorm.Open("mysql", user+":"+pass+"@/"+name+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return db, nil
}
