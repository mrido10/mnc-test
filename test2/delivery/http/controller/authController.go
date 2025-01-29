package controller

import (
	"github.com/gofiber/fiber/v2"
	"test2/model"
	"test2/usecase"
)

type AuthController struct {
	auth usecase.Auths
}

func NewAuthController(uc usecase.Auths) AuthController {
	return AuthController{auth: uc}
}

func (a AuthController) Register(c *fiber.Ctx) (interface{}, *model.Error) {
	var req model.UserRequest
	if err := c.BodyParser(&req); err != nil {
		return model.Response{}, model.NewError(400, "Body format error", nil)
	}
	return nil, a.auth.Register(req)
}

func (a AuthController) Login(c *fiber.Ctx) (interface{}, *model.Error) {
	var req model.UserRequest
	if err := c.BodyParser(&req); err != nil {
		return model.Response{}, model.NewError(400, "Body format error", nil)
	}
	return a.auth.Login(req)
}
