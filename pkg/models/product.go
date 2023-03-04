package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/bouroo/neversitup-backend-test/pkg/utils"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type TbProduct struct {
	ProductId   uint64         `json:"product_id" gorm:"primaryKey;autoIncrement:false"`
	Title       string         `json:"title" gorm:"size:127;index"`
	Description string         `json:"description" gorm:"type:text"`
	Price       float64        `json:"price" gorm:"check:price >= 0"`
	CreatedAt   sql.NullTime   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   sql.NullTime   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (t *TbProduct) BeforeCreate(tx *gorm.DB) (err error) {
	// init empty Id with snowflake Id
	if t.ProductId == 0 {
		t.ProductId = utils.SnowflakeId()
	}
	return
}

type PostProduct struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,number,gte=0"`
}

func (c *PostProduct) ValidateStruct(validate *validator.Validate) (errs []error) {
	err := validate.Struct(c)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, errors.New(fmt.Sprintf("FailedField: %s, Tag: %s, Value: %s", err.StructNamespace(), err.Tag(), err.Param())))
		}
	}
	return
}

func (c *PostProduct) ToTbProduct() (product TbProduct) {
	product = TbProduct{
		Title:       c.Title,
		Description: c.Description,
		Price:       c.Price,
	}
	return
}

type GetProduct struct {
	ProductId uint64   `params:"product_id" query:"product_id"`
	Title     string   `query:"title"`
	Available *bool    `query:"available"`
	Price     float64  `query:"price"`
	OrderBy   []string `query:"order_by"`
	Page      int      `query:"page"`
	Perpage   int      `query:"per_page"`
	Offset    int      `query:"-"`
}

func (c *GetProduct) ToTbProduct() (product TbProduct) {
	product = TbProduct{
		ProductId: c.ProductId,
		Title:     c.Title,
		Price:     c.Price,
	}
	return
}

type PresentProduct struct {
	Success    bool        `json:"success"`
	Data       []TbProduct `json:"data"`
	Errors     []Error     `json:"errors,omitempty"`
	ResultInfo *ResultInfo `json:"result_info,omitempty"`
}
