package present

import (
	"log"
	"net/http"

	"github.com/bouroo/neversitup-backend-test/pkg/models"
	"github.com/gofiber/fiber/v2"
)

var CustomErrorHandler = func(c *fiber.Ctx, err error) error {
	// Default 500 statuscode
	code := http.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		// Override status code if fiber.Error type
		code = e.Code
	}
	responseForm := models.Result{
		Errors: []models.Error{
			{
				Code: code,
				// Source:  logger.GetStackTrace(),
				Title:   http.StatusText(code),
				Message: err.Error(),
			},
		},
	}

	// Return statuscode with error message
	err = c.Status(code).JSON(responseForm)
	if err != nil {
		// In case the JSON fails
		log.Printf("customErrorHandler: %+v", err)
		return c.Status(500).SendString("Internal Server Error")
	}

	// Return from handler
	return nil
}
