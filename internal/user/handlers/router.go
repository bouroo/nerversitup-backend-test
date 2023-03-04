package handlers

import (
	"github.com/bouroo/neversitup-backend-test/internal/user/repository"
	"github.com/bouroo/neversitup-backend-test/internal/user/usecase"
	"github.com/bouroo/neversitup-backend-test/pkg/middleware"
	"github.com/bouroo/neversitup-backend-test/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, dbConn *gorm.DB) {
	userRepo := repository.NewUserRepository(dbConn)
	defaultUser := models.TbUser{
		FullName: "default user",
		Email:    viper.GetString("admin_user"),
		Password: viper.GetString("admin_password"),
		Level:    9,
	}
	userRepo.AutoMigrate(defaultUser)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandlers := NewUserHandlers(userUsecase)

	userv1 := app.Group("/api/v1/users", middleware.ReqAuth())
	{
		userv1.Get("/:email?", userHandlers.getUsers())

		userv1.Patch("/:email", userHandlers.updateUser())

		userv1.Delete("/:email", userHandlers.deleteUser())
	}
}
