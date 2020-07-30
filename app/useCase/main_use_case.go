package useCase

import (
	"app/domain/repository"
	"app/domain/service"
)

type MainUseCase struct {
	ur repository.UserRepository
	ts service.TwitterService
}

func NewMainUseCase(ur repository.UserRepository, ts service.TwitterService) *MainUseCase {
	return &MainUseCase{ur: ur}
}

func (mu MainUseCase) preLogin() {

}
