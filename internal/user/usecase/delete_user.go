package usecase

import "github.com/bouroo/neversitup-backend-test/pkg/models"

func (u *userUsecase) DeleteUser(email string) (user models.TbUser, err error) {
	user, err = u.UserRepository.DeleteUser(email)
	return
}
