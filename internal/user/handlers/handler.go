package handlers

import (
	"net/http"

	"github.com/bouroo/neversitup-backend-test/pkg/domain"
	"github.com/bouroo/neversitup-backend-test/pkg/models"
	"github.com/bouroo/neversitup-backend-test/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandlers struct {
	Validate    *validator.Validate
	UserUsecase domain.UserUsecase
}

func NewUserHandlers(userUsecase domain.UserUsecase) UserHandlers {
	return UserHandlers{
		Validate:    validator.New(),
		UserUsecase: userUsecase,
	}
}

func (h *UserHandlers) getUsers() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var reqForm models.GetUser
		var respForm models.PresentUser
		// default error code
		c.Status(http.StatusInternalServerError)
		err = c.QueryParser(&reqForm)
		if err != nil {
			err = fiber.NewError(http.StatusUnprocessableEntity, err.Error())
			return
		}

		// get email from params if assigned as param
		if email := c.Params("email"); len(email) != 0 {
			reqForm.Email = email
		}

		// force non-admin view only there profile
		if level := utils.LevelFromLocals(c.Locals("level")); level < 5 {
			reqForm.Email = c.Locals("email").(string)
		}

		users := make([]models.TbUser, 0)
		var count int
		var total int64

		// support pagination
		if reqForm.Page+reqForm.Perpage > 0 {
			if reqForm.Page == 0 {
				reqForm.Page = 1
			}

			switch {
			case reqForm.Perpage > 100:
				reqForm.Perpage = 100
			case reqForm.Perpage < 10:
				reqForm.Perpage = 10
			}

			reqForm.Offset = (reqForm.Page - 1) * reqForm.Perpage
		}

		users, count, total, err = h.UserUsecase.GetUsers(reqForm)

		respForm.Data = users
		respForm.ResultInfo = &models.ResultInfo{
			Page:      reqForm.Page,
			PerPage:   reqForm.Perpage,
			Count:     count,
			TotalCont: int(total),
		}

		if err == nil {
			c.Status(http.StatusOK)
			respForm.Success = true
		}
		return c.JSON(respForm)
	}
}

func (h *UserHandlers) updateUser() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var reqForm models.PostRegister
		var respForm models.PresentUser
		// default error code
		c.Status(http.StatusInternalServerError)

		email := c.Params("email")
		if len(email) == 0 {
			err = fiber.NewError(http.StatusBadRequest, "need email params")
			return
		}

		err = c.BodyParser(&reqForm)
		if err != nil {
			err = fiber.NewError(http.StatusUnprocessableEntity, err.Error())
			return
		}

		// force non-admin edit only there profile
		if level := utils.LevelFromLocals(c.Locals("level")); level < 5 {
			email = c.Locals("email").(string)
		}

		user, err := h.UserUsecase.UpdateUser(email, reqForm)

		respForm.Data = []models.TbUser{user}

		if err == nil {
			c.Status(http.StatusOK)
			respForm.Success = true
		}
		return c.JSON(respForm)
	}
}

func (h *UserHandlers) deleteUser() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var respForm models.PresentUser
		// default error code
		c.Status(http.StatusInternalServerError)

		email := c.Params("email")
		if len(email) == 0 {
			err = fiber.NewError(http.StatusBadRequest, "need email params")
			return
		}

		// force non-admin delete only there profile
		if level := utils.LevelFromLocals(c.Locals("level")); level < 5 {
			email = c.Locals("email").(string)
		}

		user, err := h.UserUsecase.DeleteUser(email)

		respForm.Data = []models.TbUser{user}

		if err == nil {
			c.Status(http.StatusOK)
			respForm.Success = true
		}
		return c.JSON(respForm)
	}
}
