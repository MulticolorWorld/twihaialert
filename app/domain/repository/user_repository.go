package repository

import "app/domain/entity"

type UserRepository interface {
	Create(user *entity.User) (*entity.User, error)
	Update(user *entity.User) (*entity.User, error)
	Delete(user *entity.User) error
	FindById(id int) ([]entity.User, error)
}
