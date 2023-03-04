package domain

import "github.com/bouroo/neversitup-backend-test/pkg/models"

type ProductUsecase interface {
	CreateProduct(productForms []models.PostProduct) (products []models.TbProduct, err error)
	GetProducts(condition models.GetProduct) (products []models.TbProduct, count int, total int64, err error)
	UpdateProduct(productId uint64, productForm models.PostProduct) (product models.TbProduct, err error)
	DeleteProduct(productId uint64) (product models.TbProduct, err error)
}

type ProductRepository interface {
	AutoMigrate() (err error)

	CreateProducts(products []models.TbProduct) (err error)
	GetProducts(condition models.TbProduct, orders []string, offset int, limit int) (products []models.TbProduct, count int, total int64, err error)
	UpdateProduct(productId uint64, updateValues models.TbProduct) (product models.TbProduct, err error)
	DeleteProduct(productId uint64) (product models.TbProduct, err error)
}
