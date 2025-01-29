package controller

import (
	"github.com/gofiber/fiber/v2"
	"test2/delivery/http/midleware"
)

type Route struct {
	*fiber.App
	auth AuthController
}

type RouteWithToken struct {
	*fiber.App
	transaction TransactionController
}

func NewRoute(app *fiber.App, auth AuthController) *Route {
	return &Route{
		App:  app,
		auth: auth,
	}
}

func NewRouteWithToken(app *fiber.App, transaction TransactionController) *RouteWithToken {
	return &RouteWithToken{
		App:         app,
		transaction: transaction,
	}
}

func (r *Route) SetupRoute() {
	r.Post("/register", middleware.Serve(r.auth.Register))
	r.Post("/login", middleware.Serve(r.auth.Login))
}

func (r *RouteWithToken) SetupRouteWithToken() {
	r.Post("/topup", middleware.Serve(r.transaction.TopUp))
	r.Post("/pay", middleware.Serve(r.transaction.Payment))
}
