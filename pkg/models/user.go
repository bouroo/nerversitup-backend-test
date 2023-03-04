package models

import (
	"database/sql"

	"github.com/bouroo/neversitup-backend-test/pkg/utils"
	"gorm.io/gorm"
)

type TbUser struct {
	Email     string         `json:"email" gorm:"size:64;primaryKey"`
	Password  string         `json:"password,omitempty" gorm:"-"`
	Hashed    string         `json:"hash,omitempty" gorm:"size:127"`
	FullName  string         `json:"full_name" gorm:"size:64;index"`
	Address   string         `json:"address" gorm:"type:text"`
	Level     uint           `json:"level" gorm:"size:1;default:1;index"`
	CreatedAt sql.NullTime   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt sql.NullTime   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (t *TbUser) BeforeSave(tx *gorm.DB) (err error) {
	// hashed password
	if len(t.Hashed) == 0 && len(t.Password) != 0 {
		hashed, err := utils.HashPassword(t.Password)
		if err != nil {
			return err
		}
		t.Password = ""
		t.Hashed = string(hashed)
	}
	return
}

func (t *TbUser) AfterSave(tx *gorm.DB) (err error) {
	// clear password clue
	t.Password = ""
	t.Hashed = ""
	return
}

type GetUser struct {
	Email    string   `params:"email" query:"email"`
	FullName string   `query:"full_name"`
	OrderBy  []string `query:"order_by"`
	Page     int      `query:"page"`
	Perpage  int      `query:"per_page"`
	Offset   int      `query:"-"`
}

func (c *GetUser) ToTbUser() (user TbUser, err error) {
	user = TbUser{
		Email:    c.Email,
		FullName: c.FullName,
	}
	return
}

type PresentUser struct {
	Success    bool        `json:"success"`
	Data       []TbUser    `json:"data"`
	Errors     []Error     `json:"errors,omitempty"`
	ResultInfo *ResultInfo `json:"result_info,omitempty"`
}
