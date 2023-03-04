package models

import (
	"database/sql"

	"github.com/bouroo/neversitup-backend-test/pkg/utils"
	"github.com/segmentio/encoding/json"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderPaid      OrderStatus = "paid"
	OrderShipped   OrderStatus = "shipped"
	OrderCancelled OrderStatus = "cancelled"
)

type TbOrder struct {
	OrderId    uint64         `json:"order_id" gorm:"primaryKey;autoIncrement:false"`
	Email      string         `json:"email" gorm:"size:64;index"`
	Status     string         `json:"status" gorm:"size:16;default:pending;index;check:status IN ('pending','paid','shipped','cancelled')"`
	OrderItems datatypes.JSON `json:"order_items" gorm:""`
	CreatedAt  sql.NullTime   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  sql.NullTime   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	User       *TbUser        `json:"user,omitempty" gorm:"foreignKey:Email"`
}

func (t *TbOrder) BeforeCreate(tx *gorm.DB) (err error) {
	// init empty Id with snowflake Id
	if t.OrderId == 0 {
		t.OrderId = utils.SnowflakeId()
	}
	if len(t.Status) == 0 {
		t.Status = string(OrderPending)
	}
	return
}

type OrderItem struct {
	ProductId uint64  `json:"product_id"`
	Title     string  `json:"title"`
	Amount    uint64  `json:"amount"`
	Price     float64 `json:"price"`
}

type PostOrder struct {
	Email      string      `json:"email"`
	Status     string      `json:"status"`
	OrderItems []OrderItem `json:"order_items"`
}

func (c *PostOrder) ToTbOrder(email string) (order TbOrder, err error) {
	order = TbOrder{
		Email:  email,
		Status: c.Status,
	}
	orderItems, err := json.Marshal(c.OrderItems)
	if err != nil {
		return
	}
	order.OrderItems = orderItems
	return
}

type GetOrder struct {
	OrderId uint64   `params:"order_id" query:"order_id"`
	Email   string   `query:"email"`
	Status  string   `query:"status"`
	OrderBy []string `query:"order_by"`
	Page    int      `query:"page"`
	Perpage int      `query:"per_page"`
	Offset  int      `query:"-"`
}

func (c *GetOrder) ToTbOrder() (order TbOrder, err error) {
	order = TbOrder{
		OrderId: c.OrderId,
		Email:   c.Email,
		Status:  c.Status,
	}
	return
}

type PresentOrder struct {
	Success    bool        `json:"success"`
	Data       []TbOrder   `json:"data"`
	Errors     []Error     `json:"errors,omitempty"`
	ResultInfo *ResultInfo `json:"result_info,omitempty"`
}
