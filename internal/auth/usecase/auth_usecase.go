package usecase

import (
	"github.com/bouroo/neversitup-backend-test/pkg/domain"
)

type authUsecase struct {
	domain.UserRepository
}

func NewAuthUsecase(authRepository domain.UserRepository) domain.AuthUsecase {
	return &authUsecase{
		UserRepository: authRepository,
	}
}
