package handlers

import (
	"errors"
	"net/http"

	"github.com/bouroo/neversitup-backend-test/pkg/domain"
	"github.com/bouroo/neversitup-backend-test/pkg/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProductHandlers struct {
	Validate       *validator.Validate
	ProductUsecase domain.ProductUsecase
}

func NewProductHandlers(productUsecase domain.ProductUsecase) ProductHandlers {
	return ProductHandlers{
		Validate:       validator.New(),
		ProductUsecase: productUsecase,
	}
}

func (h *ProductHandlers) createProduct() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var reqForms []models.PostProduct
		var respForm models.PresentProduct
		// default error code
		c.Status(http.StatusInternalServerError)
		err = c.BodyParser(&reqForms)
		if err != nil {
			err = fiber.NewError(http.StatusUnprocessableEntity, err.Error())
			return
		}

		for _, reqForm := range reqForms {
			if errs := reqForm.ValidateStruct(h.Validate); len(errs) != 0 {
				err = fiber.NewError(http.StatusBadRequest, errors.Join(errs...).Error())
				return
			}
		}

		product, err := h.ProductUsecase.CreateProduct(reqForms)
		if err != nil {
			return
		}

		respForm.Data = product

		if err == nil {
			c.Status(http.StatusOK)
			respForm.Success = true
		}
		return c.JSON(respForm)
	}
}

func (h *ProductHandlers) getProducts() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var reqForm models.GetProduct
		var respForm models.PresentProduct
		// default error code
		c.Status(http.StatusInternalServerError)
		err = c.QueryParser(&reqForm)
		if err != nil {
			err = fiber.NewError(http.StatusUnprocessableEntity, err.Error())
			return
		}

		// get productId from params if assigned as param
		if productId, _ := c.ParamsInt("product_id"); productId != 0 {
			reqForm.ProductId = uint64(productId)
		}

		products := make([]models.TbProduct, 0)
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

		products, count, total, err = h.ProductUsecase.GetProducts(reqForm)

		respForm.Data = products
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

func (h *ProductHandlers) updateProduct() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var reqForm models.PostProduct
		var respForm models.PresentProduct
		// default error code
		c.Status(http.StatusInternalServerError)

		productId, err := c.ParamsInt("product_id")
		if err != nil {
			err = fiber.NewError(http.StatusUnprocessableEntity, err.Error())
			return
		}

		err = c.BodyParser(&reqForm)
		if err != nil {
			err = fiber.NewError(http.StatusUnprocessableEntity, err.Error())
			return
		}

		product, err := h.ProductUsecase.UpdateProduct(uint64(productId), reqForm)

		respForm.Data = []models.TbProduct{product}

		if err == nil {
			c.Status(http.StatusOK)
			respForm.Success = true
		}
		return c.JSON(respForm)
	}
}

func (h *ProductHandlers) deleteProduct() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var respForm models.PresentProduct
		// default error code
		c.Status(http.StatusInternalServerError)
		productId, err := c.ParamsInt("product_id")
		if err != nil {
			err = fiber.NewError(http.StatusUnprocessableEntity, err.Error())
			return
		}

		product, err := h.ProductUsecase.DeleteProduct(uint64(productId))
		if err != nil {
			err = fiber.NewError(http.StatusInternalServerError, err.Error())
			return
		}

		respForm.Data = []models.TbProduct{product}

		if err == nil {
			c.Status(http.StatusOK)
			respForm.Success = true
		}
		return c.JSON(respForm)
	}
}
