package handlers

import (
	"github.com/bouroo/neversitup-backend-test/internal/auth/repository"
	"github.com/bouroo/neversitup-backend-test/internal/auth/usecase"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, dbConn *gorm.DB) {
	authRepo := repository.NewAuthRepository(dbConn)
	authUsecase := usecase.NewAuthUsecase(authRepo)
	authHandlers := NewAuthHandlers(authUsecase)

	authRoutes := app.Group("/auth")
	{
		authRoutes.Post("/register", authHandlers.register())

		authRoutes.Post("/login", authHandlers.login())
	}
}
