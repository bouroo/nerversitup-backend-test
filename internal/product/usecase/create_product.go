package usecase

import (
	"github.com/bouroo/neversitup-backend-test/pkg/models"
)

func (u *productUsecase) CreateProduct(productForms []models.PostProduct) (products []models.TbProduct, err error) {
	for _, productForm := range productForms {
		products = append(products, productForm.ToTbProduct())
	}
	err = u.ProductRepository.CreateProducts(products)
	return
}
