package handlers

import (
	"github.com/bouroo/neversitup-backend-test/internal/product/repository"
	"github.com/bouroo/neversitup-backend-test/internal/product/usecase"
	"github.com/bouroo/neversitup-backend-test/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, dbConn *gorm.DB) {
	productRepo := repository.NewProductRepository(dbConn)
	productRepo.AutoMigrate()
	productUsecase := usecase.NewProductUsecase(productRepo)
	productHandlers := NewProductHandlers(productUsecase)

	apiv1 := app.Group("/api/v1/products")
	{
		apiv1.Post("/", middleware.ReqAuth(5), productHandlers.createProduct())

		apiv1.Get("/:product_id<int>?", productHandlers.getProducts())

		apiv1.Patch("/:product_id<int>", middleware.ReqAuth(5), productHandlers.updateProduct())

		apiv1.Delete("/:product_id<int>", middleware.ReqAuth(5), productHandlers.deleteProduct())
	}
}
