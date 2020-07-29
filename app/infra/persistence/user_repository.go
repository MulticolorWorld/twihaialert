package persistence

import (
	"app/domain/entity"
	"app/domain/repository"
	"github.com/jinzhu/gorm"
	"os"
)

var connectedDB *gorm.DB = nil

type UserPersistence struct {
	db *gorm.DB
}

func NewUserPersistence() repository.UserRepository {
	return &UserPersistence{db: connectedDB}
}

func InitDBConnection() error {
	user := os.Getenv("twihaialert_db_user")
	pass := os.Getenv("twihaialert_db_password")
	name := os.Getenv("twihaialert_db_name")

	db, err := gorm.Open("mysql", user+":"+pass+"@/"+name+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return err
	}

	db.AutoMigrate(&entity.User{})
	connectedDB = db
	return nil
}

func (up UserPersistence) Create(u *entity.User) (*entity.User, error) {
	err := up.db.Create(u).Error
	return u, err
}
