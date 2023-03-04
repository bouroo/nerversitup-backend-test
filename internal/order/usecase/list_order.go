package usecase

import "github.com/bouroo/neversitup-backend-test/pkg/models"

func (u *orderUsecase) GetOrders(condition models.GetOrder) (orders []models.TbOrder, count int, total int64, err error) {

	order, err := condition.ToTbOrder()
	if err != nil {
		return
	}

	orders, count, total, err = u.OrderRepository.GetOrders(order, condition.OrderBy, condition.Offset, condition.Perpage)

	if err != nil {
		return
	}

	return
}
