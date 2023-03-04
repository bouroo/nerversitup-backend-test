package usecase

import (
	"github.com/bouroo/neversitup-backend-test/pkg/domain"
)

type userUsecase struct {
	UserRepository domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{UserRepository: userRepo}
}
