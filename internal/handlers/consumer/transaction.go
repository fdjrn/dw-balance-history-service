package consumer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/fdjrn/dw-balance-history-service/internal/db/entity"
	"github.com/fdjrn/dw-balance-history-service/internal/db/repository"
	"github.com/fdjrn/dw-balance-history-service/internal/utilities"
	"time"
)

type TransactionHandler struct {
	repository repository.BalanceHistoryRepository
}

func NewTransactionHandler() TransactionHandler {
	return TransactionHandler{repository: repository.NewBalanceHistoryRepository()}
}

func (h *TransactionHandler) DoHandleTransaction(message *sarama.ConsumerMessage) (*entity.BalanceHistory, error) {
	data := new(entity.BalanceTransaction)
	err := json.Unmarshal(message.Value, data)
	if err != nil {
		return nil, err
	}

	// cek success status "00"
	if data.Status != "00" {
		return nil, errors.New("| cannot process current transaction. unsuccessful status")
	}
	// cek duplicate insert data based on transType
	if h.repository.IsExists(data.ReceiptNumber) {
		return nil, errors.New(
			fmt.Sprintf("| transaction with receipt number %s, already exists. insert document skipped...",
				data.ReceiptNumber),
		)
	}
	// populate transaction data
	h.repository.Entity.ID = ""

	td, err := time.Parse("20060102150405", data.TransDate)
	h.repository.Entity.TransDate = td.Format("2006-01-02 15:04:05")

	h.repository.Entity.TransType = data.TransType
	h.repository.Entity.Description = data.Items[0].Name
	h.repository.Entity.PartnerID = data.PartnerID
	h.repository.Entity.MerchantID = data.MerchantID
	h.repository.Entity.TerminalID = data.TerminalID
	h.repository.Entity.TerminalName = data.TerminalName
	h.repository.Entity.PartnerRefNumber = data.PartnerRefNumber
	h.repository.Entity.ReceiptNumber = data.ReceiptNumber

	switch data.TransType {
	case utilities.TransTypeTopUp:
		h.repository.Entity.TransCode = utilities.TransCodeTopup
		h.repository.Entity.Debit = 0
		h.repository.Entity.Credit = data.TotalAmount
	case utilities.TransTypePayment:
		h.repository.Entity.TransCode = utilities.TransCodeDeduct
		h.repository.Entity.Debit = data.TotalAmount
		h.repository.Entity.Credit = 0
	default:
		return nil, errors.New("| unknown transType value. transaction cannot be processed")
	}

	h.repository.Entity.Balance = data.LastBalance

	ts := time.Now().UnixMilli()
	h.repository.Entity.CreatedAt = ts
	h.repository.Entity.UpdatedAt = ts

	insertedId, err := h.repository.Create()
	if err != nil {
		return nil, err
	}

	// fetch inserted document
	result, err := h.repository.FindByID(insertedId)
	if err != nil {
		return nil, errors.New("| cannot fetch inserted document, or its empty")
	}

	return result, nil
}
