package handlers

import (
	"github.com/fdjrn/dw-balance-history-service/internal/db/entity"
	"github.com/fdjrn/dw-balance-history-service/internal/db/repository"
	"github.com/fdjrn/dw-balance-history-service/internal/handlers/validator"
	"github.com/gofiber/fiber/v2"
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

func (b *BalanceHistoryHandler) GetBalanceHistories(c *fiber.Ctx) error {

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

//func GetHistoryByPeriod(c *fiber.Ctx) error {
//	var request = new(HistoryRequest)
//
//	// parse body payload
//	if err := c.BodyParser(&request); err != nil {
//		return c.Status(fiber.StatusBadRequest).JSON(ResponsePayload{
//			Success: false,
//			Message: err.Error(),
//			Data: ResponsePayloadData{
//				Total:  0,
//				Result: nil,
//			},
//		})
//	}
//
//	if request.UID == "" {
//		return c.Status(fiber.StatusBadRequest).JSON(ResponsePayload{
//			Success: false,
//			Message: "uniqueId cannot be empty",
//			Data: ResponsePayloadData{
//				Total:  0,
//				Result: nil,
//			},
//		})
//	}
//
//	if !isValidPeriod(request.Period) {
//		return c.Status(fiber.StatusBadRequest).JSON(ResponsePayload{
//			Success: false,
//			Message: "invalid period",
//			Data: ResponsePayloadData{
//				Total:  0,
//				Result: nil,
//			},
//		})
//	}
//
//	code, histories, length, err := repository.BalanceHistoryRepository.FindByPeriod(request)
//	if err != nil {
//		return c.Status(code).JSON(ResponsePayload{
//			Success: false,
//			Message: err.Error(),
//			Data: ResponsePayloadData{
//				Total:  length,
//				Result: histories,
//			},
//		})
//	}
//
//	return c.Status(fiber.StatusOK).JSON(ResponsePayload{
//		Success: true,
//		Message: "balance histories fetched successfully",
//		Data: ResponsePayloadData{
//			Total:  length,
//			Result: histories,
//		},
//	})
//}
