package routes

import (
	"github.com/fdjrn/dw-balance-history-service/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func initBalanceHistoryRoutes(router fiber.Router) {
	r := router.Group("/account/balance/history")

	r.Post("/last-transaction", func(c *fiber.Ctx) error {
		return handlers.GetHistoryByLastTransaction(c)
	})

	// balance transaction
	// ---------------------------------------------------------------
	//r.Post("/topup", func(c *fiber.Ctx) error {
	//	return balances.TopupBalance(c)
	//})
	//
	//r.Post("/deduct", func(c *fiber.Ctx) error {
	//	return balances.DeductBalance(c)
	//})

}
