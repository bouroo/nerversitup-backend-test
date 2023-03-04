package usecase

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bouroo/neversitup-backend-test/pkg/models"
	"github.com/bouroo/neversitup-backend-test/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

func (u *authUsecase) Login(email, password string) (tokenStr string, err error) {
	// get user from databse
	users, count, _, err := u.UserRepository.GetUsers(models.TbUser{Email: email}, nil, 0, 0)
	if err != nil {
		return
	}
	if count == 0 {
		err = fiber.NewError(http.StatusUnauthorized, "invalid email")
		return
	}
	// validate password hashed
	if err = utils.CheckPassword(password, users[0].Hashed); err != nil {
		err = fiber.NewError(http.StatusUnauthorized, "invalid password")
		return
	}

	// access_token 1h follow
	// https://learn.microsoft.com/en-us/azure/active-directory/develop/active-directory-configurable-token-lifetimes#access-id-and-saml2-token-lifetime-policy-properties
	claims := jwt.StandardClaims{
		Subject:   email,
		Audience:  strconv.Itoa(int(users[0].Level)),
		ExpiresAt: time.Now().Add(viper.GetDuration("jwt_time")).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	tokenStr, err = token.SignedString([]byte(viper.GetString("jwt_secret")))

	return
}
