package routes

import (
	"github.com/fdjrn/dw-balance-history-service/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func initBalanceHistoryRoutes(router fiber.Router) {
	historyHandler := handlers.NewBalanceHistoryHandler()
	historyRoute := router.Group("/account/balance/history")

	//r.Post("/last-transaction", func(c *fiber.Ctx) error {
	//	return handlers.GetHistoryByLastTransaction(c)
	//})

	historyRoute.Post("/all", func(c *fiber.Ctx) error {
		return historyHandler.GetBalanceHistories(c)
	})

	//r.Post("/periods", func(c *fiber.Ctx) error {
	//	return handlers.GetHistoryByPeriod(c)
	//})

}
