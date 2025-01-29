package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"test2/model"
	"test2/repository/cacheRepository"
	"test2/util"
)

func UseToken(jwtToken util.JwtToken, cache cacheRepository.Cache) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var err *model.Error
		var resp = model.Response{
			Message: "Success",
		}
		defer func() {
			if r := recover(); r != nil {
				err = model.NewError(
					500,
					"Internal Server Error",
					errors.New(fmt.Sprintf("\n%v \n%s", r, string(debug.Stack()))))
			}

			resp, _ = c.Locals("response").(model.Response)
			if err != nil {
				resp = model.Response{
					Message: err.Message,
				}
				log.Println("[ERROR]", err.Error(), err.ErrorFile)
				c.Status(err.Code)
			}
			_ = c.JSON(resp)
		}()

		reqHeader := c.GetReqHeaders()
		if len(reqHeader) == 0 {
			return fiber.ErrUnauthorized
		}

		authorization := reqHeader["Authorization"]
		if authorization == nil {
			err = model.NewError(http.StatusUnauthorized, "Unauthorized", nil)
			return nil
		}

		tokens := strings.Split(authorization[0], " ")
		if len(tokens) != 2 {
			return model.NewError(401, "Unauthorized", nil)
		}
		if tokens[0] != "Bearer" {
			return model.NewError(401, "Unauthorized", nil)
		}

		_, err = jwtToken.ParseToken(tokens[1])
		if err != nil {
			return nil
		}

		authInfo, err := cache.Get(tokens[1])
		if err != nil {
			if err.Code == 404 {
				err = model.NewError(http.StatusUnauthorized, "Unauthorized", nil)
				return nil
			}
			err = model.NewError(http.StatusInternalServerError, "Internal server error", nil)
			return nil
		}

		var userAccess model.UserAccess
		_ = json.Unmarshal([]byte(authInfo), &userAccess)
		c.Locals("auth", userAccess)
		return c.Next()
	}
}

func Serve(serve func(*fiber.Ctx) (interface{}, *model.Error)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var err *model.Error
		var resp = model.Response{
			Message: "Success",
		}
		defer func() {
			if r := recover(); r != nil {
				err = model.NewError(
					500,
					"Internal Server Error",
					errors.New(fmt.Sprintf("\n%v \n%s", r, string(debug.Stack()))))
			}

			if err != nil {
				resp = model.Response{
					Message: err.Message,
				}
				log.Println("[ERROR]", err.Error(), err.ErrorFile)
				c.Status(err.Code)
			}
			_ = c.Locals("response", resp)
			_ = c.JSON(resp)
		}()

		resp.Result, err = serve(c)
		return nil
	}
}
