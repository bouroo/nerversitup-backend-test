package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

func ReqAuth(reqLevel ...int) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		// extract token from header
		authHeader := c.Get(fiber.HeaderAuthorization)
		authHeaders := strings.Split(authHeader, " ")
		if len(authHeaders) != 2 || authHeaders[0] != "Bearer" {
			err = fiber.NewError(http.StatusUnauthorized, "need bearer authentication")
			return
		}

		// parse and validate token
		var cliams jwt.StandardClaims
		_, err = jwt.ParseWithClaims(authHeaders[1], &cliams, func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("jwt_secret")), nil
		})
		if err != nil {
			err = fiber.NewError(http.StatusUnauthorized, err.Error())
			return
		}

		// check level
		if len(reqLevel) != 0 {
			if currLevel, _ := strconv.Atoi(cliams.Audience); currLevel < reqLevel[0] {
				err = fiber.NewError(http.StatusForbidden, "insufficient rights to a resource")
				return
			}
		}

		// set user into local memory
		c.Locals("email", cliams.Subject)
		c.Locals("level", cliams.Audience)

		return c.Next()
	}
}
