package models

import (
	"errors"
	"fmt"

	"github.com/bouroo/neversitup-backend-test/pkg/utils"
	"github.com/go-playground/validator/v10"
)

type OauthError string

const (
	Oauth400 OauthError = "invalid_request"
	Oauth403 OauthError = "access_denied"
	Oauth401 OauthError = "unauthorized_client"
	Oauth501 OauthError = "unsupported_response_type"
	Oauth406 OauthError = "invalid_scope"
	Oauth500 OauthError = "server_error"
	Oauth503 OauthError = "temporarily_unavailable"
)

type PostRegister struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	Address  string `json:"address" validate:"required"`
}

func (c *PostRegister) ToTbUser() (user TbUser, err error) {
	user = TbUser{
		Email:    c.Email,
		FullName: c.FullName,
		Address:  c.Address,
	}
	if len(c.Password) != 0 {
		hashed, err := utils.HashPassword(c.Password)
		if err != nil {
			return user, err
		}
		user.Hashed = string(hashed)
	}
	return
}

func (c *PostRegister) ValidateStruct(validate *validator.Validate) (errs []error) {
	err := validate.Struct(c)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, errors.New(fmt.Sprintf("FailedField: %s, Tag: %s, Value: %s", err.StructNamespace(), err.Tag(), err.Param())))
		}
	}
	return
}

type PostLogin struct {
	Email     string `json:"email" form:"email"`
	Passsword string `json:"password" form:"password"`
}

type PresentOauth struct {
	AccessToken      string     `json:"access_token"`
	TokenType        string     `json:"token_type"`
	EpiresIn         int        `json:"expires_in"`
	Error            OauthError `json:"error,omitempty"`
	ErrorDescription string     `json:"error_description,omitempty"`
}
