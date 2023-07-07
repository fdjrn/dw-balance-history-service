package entity

type BalanceTopUp struct {
	// MongoDB ObjectID
	ID string `json:"accountId,omitempty" bson:"_id,omitempty"`

	// Unique ID yang didapat dari Client (MyDigiLearn)
	UniqueID string `json:"uniqueId,omitempty" bson:"uniqueId,omitempty"`

	// Amount of topup
	Amount int `json:"topupAmount" bson:"topupAmount,omitempty"`

	// Internal Ref Number, generated by system
	InRefNumber string `json:"inRefNumber" bson:"inRefNumber,omitempty"`

	// External Ref Number, generated by third party e.g. SoF
	ExRefNumber string `json:"exRefNumber" bson:"exRefNumber,omitempty"`

	// External Transaction Date / Success Date Time (Topup),
	// generated by third party e.g. SoF
	TransDate int `json:"transDate" bson:"transDate,omitempty"`

	// Balance after addition
	LastBalance int64 `json:"currentBalance,omitempty" bson:"lastBalance,omitempty"`

	// Encrypted last balance after addition
	LastBalanceEncrypted string `json:"-" bson:"-"`

	// Receipt Number
	ReceiptNumber string `json:"receiptNumber,omitempty" bson:"receiptNumber"`

	CreatedAt int64 `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt int64 `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`

	//SuccessDateTime int    `json:"SuccessDateTime"`
}

type BalanceDeduction struct {
	ID                   string `json:"accountId,omitempty" bson:"_id,omitempty"`
	UniqueID             string `json:"uniqueId," bson:"uniqueId"`
	Amount               int    `json:"amount"`
	MerchantID           int    `json:"merchantID,omitempty"`
	TransType            int    `json:"transType"`
	Description          string `json:"description"`
	InvoiceNumber        string `json:"invoiceNumber"`
	ReceiptNumber        string `json:"receiptNumber,omitempty"`
	LastBalance          int64  `json:"currentBalance,omitempty"`
	LastBalanceEncrypted string `json:"-"`
}

type BalanceHistory struct {
	ID               string `json:"accountId,omitempty" bson:"_id,omitempty"`
	TransDate        string `json:"transDate" bson:"transDate"`                         // YYYY-MM-DD hh:mm:ss
	TransDateNumeric int64  `json:"transDateNumeric,omitempty" bson:"transDateNumeric"` // unix time millis
	TransCode        string `json:"transCode" bson:"transCode"`
	TransType        int    `json:"transType" bson:"transType"`
	Description      string `json:"description" bson:"description"`
	PartnerID        string `json:"partnerId" bson:"partnerId"`
	MerchantID       string `json:"merchantId" bson:"merchantId"`
	TerminalID       string `json:"terminalId" bson:"terminalId"`
	TerminalName     string `json:"terminalName" bson:"terminalName"`
	PartnerRefNumber string `json:"partnerRefNumber" bson:"partnerRefNumber"`
	ReceiptNumber    string `json:"receiptNumber" bson:"receiptNumber"`
	Debit            int64  `json:"debit" bson:"debit"`
	Credit           int64  `json:"credit" bson:"credit"`
	Balance          int64  `json:"balance" bson:"balance"`
	CreatedAt        int64  `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt        int64  `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type TransactionItem struct {
	ID     string `json:"id,omitempty" bson:"_id,omitempty"`
	Code   string `json:"code,omitempty" bson:"code"`
	Name   string `json:"name" bson:"name"`
	Amount int64  `json:"amount" bson:"amount"`
	Price  int64  `json:"price,omitempty" bson:"price"`
	Qty    int    `json:"qty,omitempty" bson:"qty"`
}

type BalanceTransaction struct {
	ID                   string            `json:"id,omitempty" bson:"_id,omitempty"`
	TransDate            string            `json:"transDate,omitempty" bson:"transDate"`               // YYYY-MM-DD hh:mm:ss
	TransDateNumeric     int64             `json:"transDateNumeric,omitempty" bson:"transDateNumeric"` // unix time millis
	ReferenceNo          string            `json:"referenceNo,omitempty" bson:"referenceNo"`
	ReceiptNumber        string            `json:"receiptNumber,omitempty" bson:"receiptNumber"`
	LastBalance          int64             `json:"lastBalance,omitempty" bson:"lastBalance"`
	LastBalanceEncrypted string            `json:"-" bson:"-"`
	Status               string            `json:"status,omitempty" bson:"status"`
	TransType            int               `json:"transType,omitempty" bson:"transType"` // (1) TopUp | (2) Payment | (3) Distribution
	PartnerTransDate     string            `json:"partnerTransDate" bson:"partnerTransDate"`
	PartnerRefNumber     string            `json:"partnerRefNumber" bson:"partnerRefNumber"`
	PartnerID            string            `json:"partnerId" bson:"partnerId"`
	MerchantID           string            `json:"merchantId" bson:"merchantId"`
	TerminalID           string            `json:"terminalId" bson:"terminalId"`
	TerminalName         string            `json:"terminalName" bson:"terminalName"`
	TotalAmount          int64             `json:"totalAmount" bson:"totalAmount"`
	Items                []TransactionItem `json:"items" bson:"items"`
	CreatedAt            int64             `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt            int64             `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
