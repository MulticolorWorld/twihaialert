package repository

import "app/domain/entity"

type UserRepository interface {
	Create(user *entity.User) (*entity.User, error)
}
