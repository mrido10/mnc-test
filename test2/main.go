package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"test2/config"
	"test2/delivery/http/controller"
	middleware "test2/delivery/http/midleware"
	"test2/migrations"
	cache "test2/repository/cacheRepository"
	sql "test2/repository/sqlRepository"
	"test2/usecase"
	"test2/util"
)

func main() {
	conf := config.Conf
	db := sql.SqlConnection(conf.Database.Dsn, conf.Database.MaxIdle, conf.Database.MaxOpenConnection, conf.Database.MaxLifeTimeConnection)
	redisClient := cache.RedisConnection(conf.Cache.Address, conf.Cache.Password, conf.Cache.DB)

	if err := migrations.NewMigrate(db).Migrate(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "POST",
	}))
	app.Use(logger.New(logger.Config{
		Format:     fmt.Sprintf("${time} [%s] (${status}) ${method} ${path} \t\t${latency} \n", "DEBUG"),
		TimeFormat: "2006/01/02 15:04:05",
	}))

	userRepo := sql.NewUser(db)
	transactionRepo := sql.NewTransaction(db)
	sqlRepo := sql.NewSql(db)
	cacheRepo := cache.NewRedis(context.Background(), redisClient)

	jwtToken := util.NewJwtToken(conf.Auth.SecretKey)
	auth := usecase.NewAuth(userRepo, sqlRepo, cacheRepo, jwtToken)
	transaction := usecase.NewTransaction(userRepo, transactionRepo, sqlRepo, cacheRepo)

	controller.NewRoute(
		app,
		controller.NewAuthController(auth),
	).SetupRoute()

	app.Use(middleware.UseToken(jwtToken, cacheRepo))
	controller.NewRouteWithToken(
		app,
		controller.NewTransactionController(transaction),
	).SetupRouteWithToken()

	err := app.Listen(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
