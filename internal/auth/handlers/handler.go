package handlers

import (
	"errors"
	"net/http"

	"github.com/bouroo/neversitup-backend-test/pkg/domain"
	"github.com/bouroo/neversitup-backend-test/pkg/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type AuthHandlers struct {
	Validate    *validator.Validate
	AuthUsecase domain.AuthUsecase
}

func NewAuthHandlers(authUsecase domain.AuthUsecase) AuthHandlers {
	return AuthHandlers{
		Validate:    validator.New(),
		AuthUsecase: authUsecase,
	}
}

func (h *AuthHandlers) register() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var reqForm models.PostRegister
		var respForm models.PresentUser
		// default error code
		c.Status(http.StatusInternalServerError)
		err = c.BodyParser(&reqForm)
		if err != nil {
			err = fiber.NewError(http.StatusUnprocessableEntity, err.Error())
			return
		}

		if errs := reqForm.ValidateStruct(h.Validate); len(errs) != 0 {
			err = fiber.NewError(http.StatusBadRequest, errors.Join(errs...).Error())
			return
		}

		users, err := h.AuthUsecase.Register(reqForm)
		if err != nil {
			return
		}

		respForm.Data = users

		if err == nil {
			c.Status(http.StatusOK)
			respForm.Success = true
		}
		return c.JSON(respForm)
	}
}

func (h *AuthHandlers) login() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var reqForm models.PostLogin
		var respForm models.PresentOauth
		// default error code
		c.Status(http.StatusInternalServerError)
		respForm.Error = models.Oauth500
		err = c.BodyParser(&reqForm)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			respForm.Error = models.Oauth401
			respForm.ErrorDescription = err.Error()
			return c.JSON(respForm)
		}

		tokenStr, err := h.AuthUsecase.Login(reqForm.Email, reqForm.Passsword)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			respForm.Error = models.Oauth401
			respForm.ErrorDescription = err.Error()
			return c.JSON(respForm)
		}

		respForm.AccessToken = tokenStr
		respForm.TokenType = "Bearer"
		respForm.EpiresIn = int(viper.GetDuration("jwt_time").Seconds())

		if err == nil {
			c.Status(http.StatusOK)
			respForm.Error = ""
		}
		return c.JSON(respForm)
	}
}
