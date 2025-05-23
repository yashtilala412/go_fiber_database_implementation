package utils

import (
	"clevergo.tech/jsend"
	"github.com/gofiber/fiber/v2"
)

// JSONResponse is a generic struct to unmarshal JSend responses for testing
type JSONResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
}

// JSONSuccess is a generic success output writer
func JSONSuccess(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(jsend.New(data))
}

// JSONFail is a generic fail output writer
// JSONFail can used for 4xx status code response
func JSONFail(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(jsend.NewFail(data))
}

// JSONError is a generic error output writer
// JSONError can used for 5xx status code response
func JSONError(c *fiber.Ctx, statusCode int, err string) error {
	return c.Status(statusCode).JSON(jsend.NewError(err, statusCode, nil))
}
