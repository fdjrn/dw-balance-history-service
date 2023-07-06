package handlers

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

func isValidPeriod(r HistoryRequestPeriod) bool {
	if r.Year == 0 {
		return false
	}

	if r.Month == 0 || r.Month > 12 {
		return false
	}

	return true
}

//func InsertDeductHistory(message *sarama.ConsumerMessage) (*entity.BalanceHistory, error) {
//	data := new(entity.BalanceDeduction)
//
//	err := json.Unmarshal(message.Value, &data)
//	if err != nil {
//		log.Println(err.Error())
//		return nil, err
//	}
//
//	// check for duplicate insert
//	if repository.BalanceHistoryRepository.IsExists(data.ReceiptNumber) {
//		return nil, errors.New(
//			fmt.Sprintf("transaction with receipt number %s, already exists. insert document skipped...",
//				data.ReceiptNumber),
//		)
//	}
//
//	history := new(entity.BalanceHistory)
//	history.ID = ""
//	history.UniqueID = data.UniqueID
//	history.TransDate = time.Now().UnixMilli()
//	history.TransCode = utilities.TransCodeDeduct
//	history.Description = data.Description
//	history.MerchantID = "10000"
//	history.InvoiceNumber = data.InvoiceNumber
//	history.ReceiptNumber = data.ReceiptNumber
//	history.Debit = data.Amount
//	history.Credit = 0
//	history.Balance = data.LastBalance
//	history.CreatedAt = time.Now().UnixMilli()
//	history.UpdatedAt = time.Now().UnixMilli()
//
//	_, insertedId, err := repository.BalanceHistoryRepository.InsertBalanceHistory(history)
//	if err != nil {
//		return nil, err
//	}
//
//	// fetch inserted document
//	doc, err := repository.BalanceHistoryRepository.FindByID(insertedId)
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			return nil, errors.New("cannot fetch inserted document, or its empty")
//		}
//		return nil, err
//	}
//
//	return doc, nil
//}

//func InsertTopUpHistory(message *sarama.ConsumerMessage) (*entity.BalanceHistory, error) {
//	data := new(entity.BalanceTopUp)
//
//	err := json.Unmarshal(message.Value, &data)
//	if err != nil {
//		return nil, err
//	}
//
//	// check for duplicate insert
//	if repository.BalanceHistoryRepository.IsExists(data.ReceiptNumber) {
//		return nil, errors.New(
//			fmt.Sprintf("transaction with receipt number %s, already exists. insert document skipped...",
//				data.ReceiptNumber),
//		)
//	}
//
//	history := new(entity.BalanceHistory)
//	history.ID = ""
//	history.UniqueID = data.UniqueID
//	history.TransDate = time.Now().UnixMilli()
//	history.TransCode = utilities.TransCodeTopup
//	history.Description = "Pembelian Voucher (Topup Saldo)"
//	history.MerchantID = "10000"
//	history.InvoiceNumber = data.ExRefNumber
//	history.ReceiptNumber = data.ReceiptNumber
//	history.Debit = 0
//	history.Credit = data.Amount
//	history.Balance = data.LastBalance
//	history.CreatedAt = time.Now().UnixMilli()
//	history.UpdatedAt = time.Now().UnixMilli()
//
//	_, insertedId, err := repository.BalanceHistoryRepository.InsertBalanceHistory(history)
//	if err != nil {
//		return nil, err
//	}
//
//	// fetch inserted document
//	doc, err := repository.BalanceHistoryRepository.FindByID(insertedId)
//	if err != nil {
//		if err == mongo.ErrNoDocuments {
//			return nil, errors.New("cannot fetch inserted document, or its empty")
//		}
//		return nil, err
//	}
//
//	return doc, nil
//}

//func GetHistoryByLastTransaction(c *fiber.Ctx) error {
//
//	var request = new(HistoryRequest)
//
//	// parse body payload
//	if err := c.BodyParser(request); err != nil {
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
//	// check limits parameter
//	if !isValidLimit(request.Limit) {
//		return c.Status(400).JSON(ResponsePayload{
//			Success: false,
//			Message: "valid limit value are 5, 10, 20, 50",
//			Data: ResponsePayloadData{
//				Total:  0,
//				Result: nil,
//			},
//		})
//	}
//
//	code, histories, length, err := repository.BalanceHistoryRepository.FindByLastTransaction(request)
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
//	if length == 0 {
//		return c.Status(fiber.StatusOK).JSON(ResponsePayload{
//			Success: true,
//			Message: "no document found or its empty",
//			Data: ResponsePayloadData{
//				Total:  0,
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
//
//}

//func GetHistoryWithPagination(c *fiber.Ctx) error {
//
//	var request = new(HistoryRequestPaginated)
//
//	// parse body payload
//	if err := c.BodyParser(request); err != nil {
//		return c.Status(fiber.StatusBadRequest).JSON(ResponsePayloadPaginated{
//			Success: false,
//			Message: err.Error(),
//			Data:    ResponsePayloadDataPaginated{},
//		})
//	}
//
//	// check uniqueId parameter is not empty
//	if request.UID == "" {
//		return c.Status(fiber.StatusBadRequest).JSON(ResponsePayloadPaginated{
//			Success: false,
//			Message: "uniqueId cannot be empty",
//			Data:    ResponsePayloadDataPaginated{},
//		})
//	}
//
//	code, histories, total, pages, err := repository.BalanceHistoryRepository.FindAllPaginated(request)
//	if err != nil {
//		return c.Status(code).JSON(ResponsePayloadPaginated{
//			Success: false,
//			Message: err.Error(),
//			Data:    ResponsePayloadDataPaginated{},
//		})
//	}
//
//	return c.Status(fiber.StatusOK).JSON(ResponsePayloadPaginated{
//		Success: true,
//		Message: "balance histories fetched successfully",
//		Data: ResponsePayloadDataPaginated{
//			Result:      histories,
//			Total:       total,
//			PerPage:     request.Size,
//			CurrentPage: request.Page,
//			LastPage:    pages,
//		},
//	})
//
//}

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
