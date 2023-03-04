package usecase

import "github.com/bouroo/neversitup-backend-test/pkg/models"

func (u *productUsecase) DeleteProduct(productId uint64) (product models.TbProduct, err error) {
	product, err = u.ProductRepository.DeleteProduct(productId)
	return
}
