package repository

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bouroo/neversitup-backend-test/pkg/domain"
	"github.com/bouroo/neversitup-backend-test/pkg/logger"
	"github.com/bouroo/neversitup-backend-test/pkg/models"
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/encoding/json"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type orderRepository struct {
	DbConn     *gorm.DB
	HttpClient *fasthttp.Client
}

func NewOrderRepository(dbConn *gorm.DB) domain.OrderRepository {
	return &orderRepository{
		DbConn: dbConn,
		HttpClient: &fasthttp.Client{
			ReadTimeout: 8 * time.Second,
		},
	}
}

func (r *orderRepository) AutoMigrate() (err error) {
	err = r.DbConn.AutoMigrate(
		&models.TbOrder{},
	)
	if err != nil {
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
	}
	return
}

func (r *orderRepository) CreateOrders(orders []models.TbOrder) (err error) {
	if r.DbConn == nil {
		err = errors.New("database has gone away")
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	dbTx := r.DbConn.Begin()
	defer dbTx.Rollback()

	err = dbTx.CreateInBatches(orders, 10).Error
	if err != nil {
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	return dbTx.Commit().Error
}

func (r *orderRepository) GetOrder(orderId uint64) (tbOrder models.TbOrder, err error) {
	if r.DbConn == nil {
		err = errors.New("database has gone away")
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	dbTx := r.DbConn.Model(&models.TbOrder{})
	dbTx = dbTx.Where(models.TbOrder{OrderId: orderId})
	err = dbTx.Take(&tbOrder).Error
	if err != nil {
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	return
}

func (r *orderRepository) GetOrders(condition models.TbOrder, orders []string, offset int, limit int) (tbOrders []models.TbOrder, count int, total int64, err error) {
	if r.DbConn == nil {
		err = errors.New("database has gone away")
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	dbTx := r.DbConn.Model(&models.TbOrder{})

	if condition.OrderId != 0 {
		dbTx = dbTx.Unscoped()
		dbTx = dbTx.Where(models.TbOrder{OrderId: condition.OrderId})
	} else {
		if len(condition.Email) != 0 {
			dbTx = dbTx.Where(models.TbOrder{Email: condition.Email})
		}

		if len(condition.Status) != 0 {
			dbTx = dbTx.Where(models.TbOrder{Status: condition.Status})
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

	err = dbTx.Find(&tbOrders).Error
	if err != nil {
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	count = len(tbOrders)

	return
}

func (r *orderRepository) UpdateOrder(orderId uint64, updateValues models.TbOrder) (order models.TbOrder, err error) {
	if r.DbConn == nil {
		err = errors.New("database has gone away")
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	dbTx := r.DbConn.Begin()
	defer dbTx.Rollback()

	order.OrderId = orderId
	dbTx = dbTx.Model(&order).Clauses(clause.Returning{})

	err = dbTx.Updates(updateValues).Error
	if err != nil {
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	err = dbTx.Commit().Error

	return
}

func (r *orderRepository) CancelOrder(email string, orderId uint64) (order models.TbOrder, err error) {
	if r.DbConn == nil {
		err = errors.New("database has gone away")
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	dbTx := r.DbConn.Begin()
	defer dbTx.Rollback()

	order.OrderId = orderId
	dbTx = dbTx.Model(&order).Clauses(clause.Returning{})

	if len(email) != 0 {
		dbTx = dbTx.Where(models.TbOrder{OrderId: orderId, Email: email})
	}

	err = dbTx.Updates(models.TbOrder{Status: string(models.OrderCancelled)}).Error
	if err != nil {
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	err = dbTx.Commit().Error

	return
}

func (r *orderRepository) GetProduct(productId uint64) (tbProduct models.TbProduct, err error) {
	// init request / response pool
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()
	req.SetRequestURI(viper.GetString("product_service") + "/api/v1/products/" + strconv.Itoa(int(productId)))
	req.Header.SetMethod(fasthttp.MethodGet)
	err = r.HttpClient.Do(req, resp)
	if err != nil {
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}
	var productPresent models.PresentProduct
	err = json.Unmarshal(resp.Body(), &productPresent)
	if err != nil {
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}
	if productPresent.ResultInfo.Count == 1 {
		tbProduct = productPresent.Data[0]
	}
	if tbProduct.ProductId == 0 {
		err = fiber.NewError(http.StatusNotFound, "product_id not found")
		return
	}
	return
}
