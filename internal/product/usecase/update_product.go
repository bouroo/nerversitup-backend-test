package usecase

import "github.com/bouroo/neversitup-backend-test/pkg/models"

func (u *productUsecase) UpdateProduct(productId uint64, productForm models.PostProduct) (product models.TbProduct, err error) {
	updateValues := productForm.ToTbProduct()
	product, err = u.ProductRepository.UpdateProduct(productId, updateValues)
	return
}
