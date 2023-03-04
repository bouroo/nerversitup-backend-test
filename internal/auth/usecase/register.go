package usecase

import "github.com/bouroo/neversitup-backend-test/pkg/models"

func (u *authUsecase) Register(userForm models.PostRegister) (users []models.TbUser, err error) {
	user, err := userForm.ToTbUser()
	if err != nil {
		return
	}
	users = append(users, user)
	err = u.UserRepository.CreateUsers(users)
	return
}
