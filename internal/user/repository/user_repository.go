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

type userRepository struct {
	DbConn *gorm.DB
}

func NewUserRepository(dbConn *gorm.DB) domain.UserRepository {
	return &userRepository{DbConn: dbConn}
}

func (r *userRepository) AutoMigrate(defaultUser models.TbUser) (err error) {
	err = r.DbConn.AutoMigrate(
		&models.TbUser{},
	)
	if err != nil {
		return
	}
	err = r.createDefaultUser(defaultUser)
	return
}

func (r *userRepository) CreateUsers(users []models.TbUser) (err error) {
	if r.DbConn == nil {
		err = errors.New("database has gone away")
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	dbTx := r.DbConn.Begin()
	defer dbTx.Rollback()

	err = dbTx.CreateInBatches(users, 10).Error
	if err != nil {
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	return dbTx.Commit().Error
}

func (r *userRepository) GetUsers(condition models.TbUser, orders []string, offset int, limit int) (users []models.TbUser, count int, total int64, err error) {
	if r.DbConn == nil {
		err = errors.New("database has gone away")
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	dbTx := r.DbConn.Model(&models.TbUser{})

	if len(condition.Email) != 0 {
		dbTx = dbTx.Unscoped()
		dbTx = dbTx.Where(models.TbUser{Email: condition.Email})
	} else {
		if len(condition.FullName) != 0 {
			dbTx = dbTx.Where("title LIKE ?", "%"+condition.FullName+"%")
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

	err = dbTx.Find(&users).Error
	if err != nil {
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	count = len(users)

	return
}

func (r *userRepository) UpdateUser(email string, updateValues models.TbUser) (user models.TbUser, err error) {
	if r.DbConn == nil {
		err = errors.New("database has gone away")
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	dbTx := r.DbConn.Begin()
	defer dbTx.Rollback()

	user.Email = email
	dbTx = dbTx.Model(&user).Clauses(clause.Returning{})

	err = dbTx.Updates(updateValues).Error
	if err != nil {
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	err = dbTx.Commit().Error

	return
}

func (r *userRepository) DeleteUser(email string) (user models.TbUser, err error) {
	if r.DbConn == nil {
		err = errors.New("database has gone away")
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	dbTx := r.DbConn.Begin()
	defer dbTx.Rollback()

	user.Email = email
	err = dbTx.Clauses(clause.Returning{}).Delete(&user).Error
	if err != nil {
		log.Printf("%s\nErr: %+v", logger.WhereAmI(), err)
		return
	}

	err = dbTx.Commit().Error

	return
}

func (r *userRepository) createDefaultUser(defaultUser models.TbUser) (err error) {
	_, count, _, err := r.GetUsers(defaultUser, nil, 0, 0)
	if count != 0 {
		return
	}
	defaultUser.Level = 9
	err = r.CreateUsers([]models.TbUser{defaultUser})
	return
}
