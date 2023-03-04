package usecase

import "github.com/bouroo/neversitup-backend-test/pkg/models"

func (u *userUsecase) UpdateUser(email string, userForm models.PostRegister) (user models.TbUser, err error) {
	updateValues, err := userForm.ToTbUser()
	if err != nil {
		return
	}
	user, err = u.UserRepository.UpdateUser(email, updateValues)
	return
}
