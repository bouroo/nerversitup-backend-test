package repository

import (
	"errors"
	"log"

	"github.com/bouroo/neversitup-backend-test/pkg/domain"
	"github.com/bouroo/neversitup-backend-test/pkg/logger"
	"github.com/bouroo/neversitup-backend-test/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type productRepository struct {
	DbConn *gorm.DB
}

func NewProductRepository(dbConn *gorm.DB) domain.ProductRepository {
	return &productRepository{DbConn: dbConn}
}

func (r *productRepository) AutoMigrate() (err error) {
	err = r.DbConn.AutoMigrate(
		&models.TbProduct{},
	)
	return
}

func (r *productRepository) CreateProducts(products []models.TbProduct) (err error) {
	if r.DbConn == nil {
		err = errors.New("database has gone away")
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	dbTx := r.DbConn.Begin()
	defer dbTx.Rollback()

	err = dbTx.CreateInBatches(products, 10).Error
	if err != nil {
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	return dbTx.Commit().Error
}

func (r *productRepository) GetProduct(productId uint64) (product models.TbProduct, count int, total int64, err error) {
	if r.DbConn == nil {
		err = errors.New("database has gone away")
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	dbTx := r.DbConn.Model(&models.TbProduct{})

	dbTx = dbTx.Where(models.TbProduct{ProductId: productId})

	err = dbTx.Take(&product).Error
	if err != nil {
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	return
}

func (r *productRepository) GetProducts(condition models.TbProduct, orders []string, offset int, limit int) (products []models.TbProduct, count int, total int64, err error) {
	if r.DbConn == nil {
		err = errors.New("database has gone away")
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	dbTx := r.DbConn.Model(&models.TbProduct{})

	if condition.ProductId != 0 {
		dbTx = dbTx.Unscoped()
		dbTx = dbTx.Where(models.TbProduct{ProductId: condition.ProductId})
	} else {
		if len(condition.Title) != 0 {
			dbTx = dbTx.Where("title LIKE ?", "%"+condition.Title+"%")
		}
	}

	dbTx.Count(&total)

	for _, order := range orders {
		if len(order) != 0 {
			dbTx = dbTx.Order(order)
		}
	}

	if offset+limit != 0 {
		dbTx = dbTx.Scopes(func(db *gorm.DB) *gorm.DB {
			return db.Offset(offset).Limit(limit)
		})
	}

	err = dbTx.Find(&products).Error
	if err != nil {
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	count = len(products)

	return
}

func (r *productRepository) UpdateProduct(productId uint64, updateValues models.TbProduct) (product models.TbProduct, err error) {
	if r.DbConn == nil {
		err = errors.New("database has gone away")
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	dbTx := r.DbConn.Begin()
	defer dbTx.Rollback()

	product.ProductId = productId
	dbTx = dbTx.Model(&product).Clauses(clause.Returning{})

	err = dbTx.Updates(updateValues).Error
	if err != nil {
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	err = dbTx.Commit().Error

	return
}

func (r *productRepository) DeleteProduct(productId uint64) (product models.TbProduct, err error) {
	if r.DbConn == nil {
		err = errors.New("database has gone away")
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	dbTx := r.DbConn.Begin()
	defer dbTx.Rollback()

	product.ProductId = productId
	err = dbTx.Clauses(clause.Returning{}).Delete(&product).Error
	if err != nil {
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	err = dbTx.Commit().Error

	return
}
