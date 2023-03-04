package usecase

import (
	"github.com/bouroo/neversitup-backend-test/pkg/domain"
)

type productUsecase struct {
	ProductRepository domain.ProductRepository
}

func NewProductUsecase(productRepo domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{ProductRepository: productRepo}
}
