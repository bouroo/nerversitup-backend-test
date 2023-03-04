package handlers

import (
	"net/http"

	"github.com/bouroo/neversitup-backend-test/pkg/domain"
	"github.com/bouroo/neversitup-backend-test/pkg/models"
	"github.com/bouroo/neversitup-backend-test/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type OrderHandlers struct {
	Validate     *validator.Validate
	OrderUsecase domain.OrderUsecase
}

func NewOrderHandlers(orderUsecase domain.OrderUsecase) OrderHandlers {
	return OrderHandlers{
		Validate:     validator.New(),
		OrderUsecase: orderUsecase,
	}
}

func (h *OrderHandlers) createOrder() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var reqForm models.PostOrder
		var respForm models.PresentOrder
		// default error code
		c.Status(http.StatusInternalServerError)
		err = c.BodyParser(&reqForm)
		if err != nil {
			err = fiber.NewError(http.StatusUnprocessableEntity, err.Error())
			return
		}

		reqForm.Status = string(models.OrderPending)
		email, ok := c.Locals("email").(string)
		if !ok {
			err = fiber.NewError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		order, err := h.OrderUsecase.CreateOrder(email, reqForm)
		if err != nil {
			return
		}

		respForm.Data = order

		if err == nil {
			c.Status(http.StatusOK)
			respForm.Success = true
		}
		return c.JSON(respForm)
	}
}

func (h *OrderHandlers) getOrders() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var reqForm models.GetOrder
		var respForm models.PresentOrder
		// default error code
		c.Status(http.StatusInternalServerError)
		err = c.QueryParser(&reqForm)
		if err != nil {
			err = fiber.NewError(http.StatusUnprocessableEntity, err.Error())
			return
		}

		email, ok := c.Locals("email").(string)
		if !ok {
			err = fiber.NewError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		// force non-admin view only there order
		if level := utils.LevelFromLocals(c.Locals("level")); level < 5 {
			reqForm.Email = email
		}

		// get orderId from params if assigned as param
		if orderId, _ := c.ParamsInt("order_id"); orderId != 0 {
			reqForm.OrderId = uint64(orderId)
		}

		orders := make([]models.TbOrder, 0)
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

		orders, count, total, err = h.OrderUsecase.GetOrders(reqForm)

		respForm.Data = orders
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

func (h *OrderHandlers) updateOrder() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var reqForm models.PostOrder
		var respForm models.PresentOrder
		// default error code
		c.Status(http.StatusInternalServerError)

		err = c.BodyParser(&reqForm)
		if err != nil {
			err = fiber.NewError(http.StatusUnprocessableEntity, err.Error())
			return
		}

		orderId, err := c.ParamsInt("order_id")
		if err != nil {
			err = fiber.NewError(http.StatusUnprocessableEntity, err.Error())
			return
		}

		email, ok := c.Locals("email").(string)
		if !ok {
			err = fiber.NewError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		order, err := h.OrderUsecase.UpdateOrder(email, uint64(orderId), reqForm)

		respForm.Data = []models.TbOrder{order}

		if err == nil {
			c.Status(http.StatusOK)
			respForm.Success = true
		}
		return c.JSON(respForm)
	}
}

func (h *OrderHandlers) cancelOrder() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var respForm models.PresentOrder
		// default error code
		c.Status(http.StatusInternalServerError)
		orderId, err := c.ParamsInt("order_id")
		if err != nil {
			err = fiber.NewError(http.StatusUnprocessableEntity, err.Error())
			return
		}

		var email string
		var ok bool
		// force non-admin cancel only there order
		if level := utils.LevelFromLocals(c.Locals("level")); level < 5 {
			email, ok = c.Locals("email").(string)
			if !ok {
				err = fiber.NewError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}
		}

		order, err := h.OrderUsecase.CancelOrder(email, uint64(orderId))
		if err != nil {
			err = fiber.NewError(http.StatusInternalServerError, err.Error())
			return
		}

		respForm.Data = []models.TbOrder{order}

		if err == nil {
			c.Status(http.StatusOK)
			respForm.Success = true
		}
		return c.JSON(respForm)
	}
}
