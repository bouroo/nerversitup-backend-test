package usecase

import "github.com/bouroo/neversitup-backend-test/pkg/models"

func (u *productUsecase) GetProducts(condition models.GetProduct) (products []models.TbProduct, count int, total int64, err error) {

	products, count, total, err = u.ProductRepository.GetProducts(condition.ToTbProduct(), condition.OrderBy, condition.Offset, condition.Perpage)

	if err != nil {
		return
	}

	return
}
