package domain

import "github.com/bouroo/neversitup-backend-test/pkg/models"

type AuthUsecase interface {
	Register(userForm models.PostRegister) (users []models.TbUser, err error)
	Login(email, password string) (tokenStr string, err error)
}

type AuthRepository interface{}
