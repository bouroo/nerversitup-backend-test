package usecase

import "github.com/bouroo/neversitup-backend-test/pkg/models"

func (u *userUsecase) GetUsers(condition models.GetUser) (users []models.TbUser, count int, total int64, err error) {

	user, err := condition.ToTbUser()
	if err != nil {
		return
	}

	users, count, total, err = u.UserRepository.GetUsers(user, condition.OrderBy, condition.Offset, condition.Perpage)

	if err != nil {
		return
	}

	return
}
