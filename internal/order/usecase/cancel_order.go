package usecase

import "github.com/bouroo/neversitup-backend-test/pkg/models"

func (u *orderUsecase) CancelOrder(email string, orderId uint64) (order models.TbOrder, err error) {
	order, err = u.OrderRepository.CancelOrder(email, orderId)
	return
}
