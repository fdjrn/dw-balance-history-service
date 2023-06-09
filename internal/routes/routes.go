package routes

import (
	"errors"
	"fmt"
	"github.com/fdjrn/dw-balance-history-service/configs"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {

	api := app.Group("/api/v1")
	initBalanceHistoryRoutes(api)
}

func Initialize() error {
	config := configs.MainConfig.APIServer

	app := fiber.New(
	//fiber.Config{DisableStartupMessage: true},
	)

	setupRoutes(app)

	err := app.Listen(fmt.Sprintf(":%s", config.Port))
	if err != nil {
		return errors.New(fmt.Sprintf("error on starting service: %s", err.Error()))
	}

	return nil

}
