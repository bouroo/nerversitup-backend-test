package handlers

import (
	"github.com/bouroo/neversitup-backend-test/internal/order/repository"
	"github.com/bouroo/neversitup-backend-test/internal/order/usecase"
	"github.com/bouroo/neversitup-backend-test/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, dbConn *gorm.DB) {
	orderRepo := repository.NewOrderRepository(dbConn)
	orderRepo.AutoMigrate()
	orderUsecase := usecase.NewOrderUsecase(orderRepo)
	orderHandlers := NewOrderHandlers(orderUsecase)

	apiv1 := app.Group("/api/v1/orders")
	{
		apiv1.Post("/", middleware.ReqAuth(1), orderHandlers.createOrder())

		apiv1.Get("/:order_id<int>?", middleware.ReqAuth(1), orderHandlers.getOrders())

		apiv1.Patch("/:order_id<int>", middleware.ReqAuth(1), orderHandlers.updateOrder())

		apiv1.Delete("/:order_id<int>", middleware.ReqAuth(1), orderHandlers.cancelOrder())
	}

}
