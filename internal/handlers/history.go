package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/fdjrn/dw-balance-history-service/internal/db/entity"
	"github.com/fdjrn/dw-balance-history-service/internal/db/repository"
	"github.com/fdjrn/dw-balance-history-service/pkg/payload"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func isValidLimit(i int64) bool {
	// 5, 10, 20, 50
	limits := []int64{5, 10, 20, 50}
	for _, v := range limits {
		if i == v {
			return true
		}
	}

	return false
}

func isValidPeriod(r payload.HistoryRequestPeriod) bool {
	if r.Year == 0 {
		return false
	}

	if r.Month == 0 || r.Month > 12 {
		return false
	}

	return true
}

func InsertDeductHistory(message *sarama.ConsumerMessage) (*entity.BalanceHistory, error) {
	data := new(entity.BalanceDeduction)

	err := json.Unmarshal(message.Value, &data)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	// check for duplicate insert
	if repository.BalanceHistoryRepository.IsExists(data.ReceiptNumber) {
		return nil, errors.New(
			fmt.Sprintf("transaction with receipt number %s, already exists. insert document skipped...",
				data.ReceiptNumber),
		)
	}

	history := new(entity.BalanceHistory)
	history.ID = ""
	history.UniqueID = data.UniqueID
	history.TransDate = time.Now().UnixMilli()
	history.TransCode = entity.TransCodeDeduct
	history.Description = data.Description
	history.MerchantID = "10000"
	history.InvoiceNumber = data.InvoiceNumber
	history.ReceiptNumber = data.ReceiptNumber
	history.Debit = 0
	history.Credit = data.Amount
	history.Balance = data.LastBalance
	history.CreatedAt = time.Now().UnixMilli()
	history.UpdatedAt = time.Now().UnixMilli()

	_, insertedId, err := repository.BalanceHistoryRepository.InsertBalanceHistory(history)
	if err != nil {
		return nil, err
	}

	// fetch inserted document
	doc, err := repository.BalanceHistoryRepository.FindByID(insertedId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("cannot fetch inserted document, or its empty")
		}
		return nil, err
	}

	return doc, nil
}

func InsertTopUpHistory(message *sarama.ConsumerMessage) (*entity.BalanceHistory, error) {
	data := new(entity.BalanceTopUp)

	err := json.Unmarshal(message.Value, &data)
	if err != nil {
		return nil, err
	}

	// check for duplicate insert
	if repository.BalanceHistoryRepository.IsExists(data.ReceiptNumber) {
		return nil, errors.New(
			fmt.Sprintf("transaction with receipt number %s, already exists. insert document skipped...",
				data.ReceiptNumber),
		)
	}

	history := new(entity.BalanceHistory)
	history.ID = ""
	history.UniqueID = data.UniqueID
	history.TransDate = time.Now().UnixMilli()
	history.TransCode = entity.TransCodeTopup
	history.Description = "Pembelian Voucher (Topup Saldo)"
	history.MerchantID = "10000"
	history.InvoiceNumber = data.ExRefNumber
	history.ReceiptNumber = data.ReceiptNumber
	history.Debit = data.Amount
	history.Credit = 0
	history.Balance = data.LastBalance
	history.CreatedAt = time.Now().UnixMilli()
	history.UpdatedAt = time.Now().UnixMilli()

	_, insertedId, err := repository.BalanceHistoryRepository.InsertBalanceHistory(history)
	if err != nil {
		return nil, err
	}

	// fetch inserted document
	doc, err := repository.BalanceHistoryRepository.FindByID(insertedId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("cannot fetch inserted document, or its empty")
		}
		return nil, err
	}

	return doc, nil
}

func GetHistoryByLastTransaction(c *fiber.Ctx) error {

	var request = new(payload.HistoryRequest)

	// parse body payload
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponsePayload{
			Success: false,
			Message: err.Error(),
			Data: ResponsePayloadData{
				Total:  0,
				Result: nil,
			},
		})
	}

	// check limits parameter
	if !isValidLimit(request.Limit) {
		return c.Status(400).JSON(ResponsePayload{
			Success: false,
			Message: "valid limit value are 5, 10, 20, 50",
			Data: ResponsePayloadData{
				Total:  0,
				Result: nil,
			},
		})
	}

	code, histories, length, err := repository.BalanceHistoryRepository.FindByLastTransaction(request)
	if err != nil {
		return c.Status(code).JSON(ResponsePayload{
			Success: false,
			Message: err.Error(),
			Data: ResponsePayloadData{
				Total:  length,
				Result: histories,
			},
		})
	}

	if length == 0 {
		return c.Status(fiber.StatusOK).JSON(ResponsePayload{
			Success: true,
			Message: "no document found or its empty",
			Data: ResponsePayloadData{
				Total:  0,
				Result: histories,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(ResponsePayload{
		Success: true,
		Message: "balance histories fetched successfully",
		Data: ResponsePayloadData{
			Total:  length,
			Result: histories,
		},
	})

}

func GetHistoryByPeriod(c *fiber.Ctx) error {
	var request = new(payload.HistoryRequest)

	// parse body payload
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponsePayload{
			Success: false,
			Message: err.Error(),
			Data: ResponsePayloadData{
				Total:  0,
				Result: nil,
			},
		})
	}

	if request.UID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ResponsePayload{
			Success: false,
			Message: "uniqueId cannot be empty",
			Data: ResponsePayloadData{
				Total:  0,
				Result: nil,
			},
		})
	}

	if !isValidPeriod(request.Period) {
		return c.Status(fiber.StatusBadRequest).JSON(ResponsePayload{
			Success: false,
			Message: "invalid period",
			Data: ResponsePayloadData{
				Total:  0,
				Result: nil,
			},
		})
	}

	code, histories, length, err := repository.BalanceHistoryRepository.FindByPeriod(request)
	if err != nil {
		return c.Status(code).JSON(ResponsePayload{
			Success: false,
			Message: err.Error(),
			Data: ResponsePayloadData{
				Total:  length,
				Result: histories,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(ResponsePayload{
		Success: true,
		Message: "balance histories fetched successfully",
		Data: ResponsePayloadData{
			Total:  length,
			Result: histories,
		},
	})
}

//
//func Register(c *fiber.Ctx) error {
//
//	// new account struct
//	a := new(entity.AccountBalance)
//
//	// parse body payload
//	if err := c.BodyParser(a); err != nil {
//		return c.Status(400).JSON(handlers.Responses{
//			Success: false,
//			Message: err.Error(),
//			Data:    nil,
//		})
//	}
//
//	if isRegistered(a.UniqueID) {
//		return c.Status(400).JSON(handlers.Responses{
//			Success: false,
//			Message: "uniqueId has already been registered",
//			Data:    a,
//		})
//	}
//
//	// set default value for accountBalance document
//	key, _ := tools.GenerateSecretKey()
//	encryptedBalance, _ := tools.Encrypt([]byte(key), fmt.Sprintf("%016s", "0"))
//
//	a.ID = ""
//	a.SecretKey = key
//	a.Active = true
//	a.LastBalance = encryptedBalance
//	a.MainAccountID = "-"
//	a.CreatedAt = time.Now().UnixMilli()
//	a.UpdatedAt = a.CreatedAt
//
//	code, id, err := r.InsertDocument(a)
//	if err != nil {
//		return c.Status(code).JSON(handlers.Responses{
//			Success: false,
//			Message: err.Error(),
//			Data:    nil,
//		})
//	}
//
//	_, createdAccount, _ := r.FindByID(id, true)
//
//	return c.Status(code).JSON(handlers.Responses{
//		Success: true,
//		Message: "account has been successfully registered",
//		Data:    createdAccount,
//	})
//}

/*


// isRegistered is a private function that check whether account id has been registered based on phoneNumber
func isRegistered(uniqueId string) bool {
	_, _, err := r.FindByUniqueID(uniqueId, true)
	if err != nil {
		// no document found, its mean it can be registered
		if err == mongo.ErrNoDocuments {
			return false
		}

		// TODO handling unknown error
		return true
	}
	return true
}

// isUnregistered is a private function that check whether account id has been unregistered or not
func isUnregistered(uniqueId string) bool {
	_, _, err := r.FindByActiveStatus(uniqueId, false)
	if err != nil {
		// no document found, its mean it can be unregistered
		if err == mongo.ErrNoDocuments {
			return false
		}

		// TODO handling unknown error
		return true
	}
	return true
}

// Register is a function that used to insert new document into collection and set active status to true.
func Register(c *fiber.Ctx) error {

	// new account struct
	a := new(entity.AccountBalance)

	// parse body payload
	if err := c.BodyParser(a); err != nil {
		return c.Status(400).JSON(handlers.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	if isRegistered(a.UniqueID) {
		return c.Status(400).JSON(handlers.Responses{
			Success: false,
			Message: "uniqueId has already been registered",
			Data:    a,
		})
	}

	// set default value for accountBalance document
	key, _ := tools.GenerateSecretKey()
	encryptedBalance, _ := tools.Encrypt([]byte(key), fmt.Sprintf("%016s", "0"))

	a.ID = ""
	a.SecretKey = key
	a.Active = true
	a.LastBalance = encryptedBalance
	a.MainAccountID = "-"
	a.CreatedAt = time.Now().UnixMilli()
	a.UpdatedAt = a.CreatedAt

	code, id, err := r.InsertDocument(a)
	if err != nil {
		return c.Status(code).JSON(handlers.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	_, createdAccount, _ := r.FindByID(id, true)

	return c.Status(code).JSON(handlers.Responses{
		Success: true,
		Message: "account has been successfully registered",
		Data:    createdAccount,
	})
}

// Unregister is a function that used to change active status to false (unregistered)
func Unregister(c *fiber.Ctx) error {

	// new u struct
	u := new(entity.UnregisterAccount)

	// parse body payload
	if err := c.BodyParser(u); err != nil {
		return c.Status(400).JSON(handlers.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// check if already been unregistered
	if isUnregistered(u.UniqueID) {
		return c.Status(400).JSON(handlers.Responses{
			Success: false,
			Message: "account has already been unregistered",
			Data:    u,
		})
	}
	code, err := r.DeactivateAccount(u)
	if err != nil {
		return c.Status(code).JSON(handlers.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	// insert into accountDeactivated collection
	code, _, err = r.InsertDeactivatedAccount(u)

	_, updatedAccount, _ := r.FindByUniqueID(u.UniqueID, false)

	return c.Status(code).JSON(handlers.Responses{
		Success: true,
		Message: "account has been successfully unregistered",
		Data:    updatedAccount,
	})
}

// Reregister is a function that used to re-activation account balance by changing active status to true
// and delete accountDeactivated collection by uniqueId.
func Reregister(c *fiber.Ctx) error {
	// TODO
	// 1. change active status to True
	// 2. remove document on accountDeactivated collection

	return nil
}

// GetAllRegisteredAccount is used to find all registered account and can be filtered with their active status
func GetAllRegisteredAccount(c *fiber.Ctx) error {
	accountStatus := ""
	queryParams := c.Query("active")
	if queryParams != "" {

		switch strings.ToLower(queryParams) {
		case "true":
			accountStatus = "active "
		case "false":
			accountStatus = "unregistered "
		default:
			return c.Status(fiber.StatusBadRequest).JSON(handlers.Responses{
				Success: false,
				Message: "invalid query param value, expected value is true or false",
				Data:    nil,
			})
		}
	}

	code, accounts, err := r.FindAll(queryParams)
	if err != nil {
		return c.Status(code).JSON(handlers.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	msgResponse := fmt.Sprintf("%saccounts fetched successfully ", accountStatus)

	return c.Status(code).JSON(handlers.Responses{
		Success: true,
		Message: msgResponse,
		Data:    accounts,
	})
}

// GetRegisteredAccount is used to find registered account with active status = true
func GetRegisteredAccount(c *fiber.Ctx) error {
	id, _ := primitive.ObjectIDFromHex(c.Params("id"))
	code, account, err := r.FindByID(id, true)

	if err != nil {
		return c.Status(code).JSON(handlers.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.Status(code).JSON(handlers.Responses{
		Success: true,
		Message: "account fetched successfully ",
		Data:    account,
	})

}

func GetRegisteredAccountByUID(c *fiber.Ctx) error {
	code, account, err := r.FindByUniqueID(c.Params("uid"), true)

	if err != nil {
		return c.Status(code).JSON(handlers.Responses{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.Status(code).JSON(handlers.Responses{
		Success: true,
		Message: "account fetched successfully ",
		Data:    account,
	})

}

*/
