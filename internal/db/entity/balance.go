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
	ID            string `json:"accountId,omitempty" bson:"_id,omitempty"`
	UniqueID      string `json:"uniqueId," bson:"uniqueId,omitempty"`
	TransDate     int64  `json:"transDate" bson:"transDate"`
	TransCode     string `json:"transCode" bson:"transCode"`
	Description   string `json:"description" bson:"description"`
	MerchantID    string `json:"merchantID" bson:"merchantId"`
	InvoiceNumber string `json:"invoiceNumber" bson:"invoiceNumber"`
	ReceiptNumber string `json:"receiptNumber" bson:"receiptNumber"`
	Debit         int    `json:"debit" bson:"debit"`
	Credit        int    `json:"credit" bson:"credit"`
	Balance       int64  `json:"balance" bson:"balance"`
	CreatedAt     int64  `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt     int64  `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

const (
	TransCodeTopup  = "1000"
	TransCodeDeduct = "2000"
)
