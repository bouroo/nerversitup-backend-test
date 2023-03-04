package domain

import "github.com/bouroo/neversitup-backend-test/pkg/models"

type OrderUsecase interface {
	CreateOrder(email string, orderForm models.PostOrder) (orders []models.TbOrder, err error)
	GetOrders(condition models.GetOrder) (orders []models.TbOrder, count int, total int64, err error)
	UpdateOrder(email string, orderId uint64, orderForm models.PostOrder) (order models.TbOrder, err error)
	CancelOrder(email string, orderId uint64) (order models.TbOrder, err error)
}

type OrderRepository interface {
	AutoMigrate() (err error)

	CreateOrders(orders []models.TbOrder) (err error)
	GetOrder(orderId uint64) (tbOrder models.TbOrder, err error)
	GetOrders(condition models.TbOrder, orders []string, offset int, limit int) (tbOrders []models.TbOrder, count int, total int64, err error)
	UpdateOrder(orderId uint64, updateValues models.TbOrder) (order models.TbOrder, err error)
	CancelOrder(email string, orderId uint64) (order models.TbOrder, err error)

	GetProduct(productId uint64) (tbProduct models.TbProduct, err error)
}
