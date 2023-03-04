package usecase

import (
	"net/http"

	"github.com/bouroo/neversitup-backend-test/pkg/models"
	"github.com/gofiber/fiber/v2"
)

func (u *orderUsecase) UpdateOrder(email string, orderId uint64, orderForm models.PostOrder) (order models.TbOrder, err error) {

	var totalPrice float64
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
			totalPrice += float64(orderItem.Amount) * product.Price
		}
	}

	if len(orderItems) == 0 {
		err = fiber.NewError(http.StatusUnprocessableEntity, "empty order_items")
		return
	}

	orderForm.TotalPrice = totalPrice
	orderForm.OrderItems = orderItems

	updateValues, err := orderForm.ToTbOrder(email)
	if err != nil {
		return
	}
	order, err = u.OrderRepository.UpdateOrder(orderId, updateValues)
	return
}
