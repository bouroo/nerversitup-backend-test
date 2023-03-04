package repository

import (
	user_repo "github.com/bouroo/neversitup-backend-test/internal/user/repository"
	"github.com/bouroo/neversitup-backend-test/pkg/domain"
	"gorm.io/gorm"
)

// reuse UserRepository to access user database
func NewAuthRepository(dbConn *gorm.DB) domain.UserRepository {
	return user_repo.NewUserRepository(dbConn)
}
