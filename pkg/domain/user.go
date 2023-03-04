package domain

import "github.com/bouroo/neversitup-backend-test/pkg/models"

type UserUsecase interface {
	GetUsers(condition models.GetUser) (users []models.TbUser, count int, total int64, err error)
	UpdateUser(email string, userForm models.PostRegister) (user models.TbUser, err error)
	DeleteUser(email string) (user models.TbUser, err error)
}

type UserRepository interface {
	AutoMigrate(defaultUser models.TbUser) (err error)

	CreateUsers(users []models.TbUser) (err error)
	GetUsers(condition models.TbUser, orders []string, offset int, limit int) (users []models.TbUser, count int, total int64, err error)
	UpdateUser(email string, updateValues models.TbUser) (user models.TbUser, err error)
	DeleteUser(email string) (user models.TbUser, err error)
}
