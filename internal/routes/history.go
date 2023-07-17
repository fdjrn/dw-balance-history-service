package routes

import (
	"github.com/fdjrn/dw-balance-history-service/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func initBalanceHistoryRoutes(router fiber.Router) {
	historyHandler := handlers.NewBalanceHistoryHandler()
	historyRoute := router.Group("/account/transaction/history")

	historyRoute.Post("/last-transaction", func(c *fiber.Ctx) error {
		return historyHandler.GetLastTransaction(c)
	})

	historyRoute.Post("/all", func(c *fiber.Ctx) error {
		return historyHandler.GetBalanceHistories(c, false)
	})

	historyRoute.Post("/periods", func(c *fiber.Ctx) error {
		return historyHandler.GetBalanceHistories(c, true)
	})

}
