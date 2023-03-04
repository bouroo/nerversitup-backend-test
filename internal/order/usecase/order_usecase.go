package usecase

import (
	"github.com/bouroo/neversitup-backend-test/pkg/domain"
)

type orderUsecase struct {
	OrderRepository domain.OrderRepository
}

func NewOrderUsecase(orderRepository domain.OrderRepository) domain.OrderUsecase {
	return &orderUsecase{OrderRepository: orderRepository}
}
