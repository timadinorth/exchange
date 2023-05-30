package util

import (
	"github.com/gofiber/fiber/v2"
)

func NewError(ctx *fiber.Ctx, status int, err error) error {
	er := HTTPError{
		Error: err.Error(),
	}
	return ctx.Status(status).JSON(er)
}

func NewErrorStr(ctx *fiber.Ctx, status int, err string) error {
	er := HTTPError{
		Error: err,
	}
	return ctx.Status(status).JSON(er)
}

type HTTPError struct {
	Error string `json:"error" example:"status bad request"`
}
