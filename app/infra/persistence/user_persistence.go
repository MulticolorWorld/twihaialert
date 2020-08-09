package persistence

import (
	"app/domain/entity"
	"app/domain/repository"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

type UserPersistence struct {
	db *gorm.DB
}

func (up UserPersistence) FindAll() ([]entity.User, error) {
	var users []entity.User
	err := up.db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (up UserPersistence) Delete(user *entity.User) error {
	err := up.db.Delete(user).Error
	if err != nil {
		return err
	}
	return nil
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
	host := os.Getenv("twihaialert_db_host")
	user := os.Getenv("twihaialert_db_user")
	pass := os.Getenv("twihaialert_db_password")
	name := os.Getenv("twihaialert_db_name")

	db, err := gorm.Open("mysql", user+":"+pass+"@("+host+")/"+name+"?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(5)

	return db, nil
}
