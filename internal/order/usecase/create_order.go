package usecase

import (
	"net/http"

	"github.com/bouroo/neversitup-backend-test/pkg/models"
	"github.com/gofiber/fiber/v2"
)

func (u *orderUsecase) CreateOrder(email string, orderForm models.PostOrder) (orders []models.TbOrder, err error) {

	orderItems := make([]models.OrderItem, 0)
	// recheck product in order_items
	for _, orderItem := range orderForm.OrderItems {
		// ignore empty items
		if orderItem.ProductId == 0 || orderItem.Amount == 0 {
			continue
		}
		if product, err := u.OrderRepository.GetProduct(orderItem.ProductId); err == nil {
			orderItems = append(orderItems, models.OrderItem{
				ProductId: product.ProductId,
				Title:     product.Title,
				Amount:    orderItem.Amount,
				Price:     product.Price,
			})
		}
	}

	if len(orderItems) == 0 {
		err = fiber.NewError(http.StatusUnprocessableEntity, "empty order_items")
		return
	}

	orderForm.OrderItems = orderItems

	order, err := orderForm.ToTbOrder(email)
	if err != nil {
		return
	}
	orders = append(orders, order)
	err = u.OrderRepository.CreateOrders(orders)
	return
}
