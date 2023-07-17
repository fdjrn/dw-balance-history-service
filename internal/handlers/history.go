package handlers

import (
	"fmt"
	"github.com/fdjrn/dw-balance-history-service/internal/db/entity"
	"github.com/fdjrn/dw-balance-history-service/internal/db/repository"
	"github.com/fdjrn/dw-balance-history-service/internal/handlers/validator"
	"github.com/gofiber/fiber/v2"
	"time"
)

type BalanceHistoryHandler struct {
	repository repository.BalanceHistoryRepository
}

func NewBalanceHistoryHandler() BalanceHistoryHandler {
	return BalanceHistoryHandler{repository: repository.NewBalanceHistoryRepository()}
}

func (b *BalanceHistoryHandler) isValidLimit(i int64) bool {
	// 5, 10, 20, 50
	limits := []int64{5, 10, 20, 50}
	for _, v := range limits {
		if i == v {
			return true
		}
	}

	return false
}

func (b *BalanceHistoryHandler) GetBalanceHistories(c *fiber.Ctx, isPeriod bool) error {

	var payload = new(entity.PaginatedRequest)

	// parse body payload
	if err := c.BodyParser(payload); err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	validation, err := validator.ValidateRequest(payload)
	if err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data: map[string]interface{}{
				"errors": validation,
			},
		})
	}

	// set default param value
	if payload.Page == 0 {
		payload.Page = 1
	}

	if payload.Size == 0 {
		payload.Size = 10
	}

	// for periods endpoint
	if isPeriod {
		if payload.Periods == nil {
			return c.Status(400).JSON(entity.Responses{
				Success: false,
				Message: "periods attributes cannot be empty",
				Data:    nil,
			})
		}

		payload.Periods.StartDate, err = time.ParseInLocation(
			"20060102150405",
			fmt.Sprintf("%s%s", payload.Periods.Start, "000000"),
			time.Now().Location(),
		)
		if err != nil {
			return c.Status(400).JSON(entity.Responses{
				Success: false,
				Message: "invalid start periods",
				Data:    nil,
			})
		}

		payload.Periods.EndDate, err = time.ParseInLocation(
			"20060102150405",
			fmt.Sprintf("%s%s", payload.Periods.End, "235959"),
			time.Now().Location(),
		)

		if err != nil {
			return c.Status(400).JSON(entity.Responses{
				Success: false,
				Message: "invalid end periods",
				Data:    nil,
			})
		}

		if payload.Periods.EndDate.Before(payload.Periods.StartDate) {
			return c.Status(400).JSON(entity.Responses{
				Success: false,
				Message: "end period cannot be less than start period",
				Data:    nil,
			})
		}
	}

	b.repository.Request = payload
	histories, total, pages, err := b.repository.FindAllPaginated()
	if err != nil {
		return c.Status(500).JSON(entity.PaginatedResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.PaginatedResponse{
		Success: true,
		Message: "transaction histories fetched successfully",
		Data: &entity.PaginatedDetailResponse{
			Result: histories,
			Total:  total,
			Pagination: &entity.PaginationInfo{
				PerPage:     payload.Size,
				CurrentPage: payload.Page,
				LastPage:    pages,
			},
		},
	})

}

func (b *BalanceHistoryHandler) GetLastTransaction(c *fiber.Ctx) error {
	var payload = new(entity.PaginatedRequest)

	// parse body payload
	if err := c.BodyParser(payload); err != nil {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// check limits parameter
	if !b.isValidLimit(payload.Limit) {
		return c.Status(400).JSON(entity.Responses{
			Success: false,
			Message: "valid limit value are 5, 10, 20, 50",
			Data:    nil,
		})
	}

	b.repository.Request = payload
	histories, totalDocs, err := b.repository.FindByLastTransaction()
	if err != nil {
		return c.Status(500).JSON(entity.PaginatedResponse{
			Success: false,
			Message: err.Error(),
			//Total:   0,
			Data: nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(entity.PaginatedResponse{
		Success: true,
		Message: "last transaction fetched successfully",
		//Total:   totalDocs,
		Data: &entity.PaginatedDetailResponse{
			Total:      int64(totalDocs),
			Result:     histories,
			Pagination: nil,
		},
	})
}
