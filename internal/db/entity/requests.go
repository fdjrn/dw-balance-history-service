package entity

import "time"

//type HistoryRequest struct {
//	UID    string               `json:"uniqueId"`
//	Limit  int64                `json:"limit,omitempty"` // 5, 10, 20, 50
//	Period HistoryRequestPeriod `json:"period,omitempty"`
//}
//
//type HistoryRequestPeriod struct {
//	Year  int `json:"year"`
//	Month int `json:"month"`
//}

type PeriodsRequest struct {
	Start     string    `json:"start,omitempty"`
	StartDate time.Time `json:"-"`
	End       string    `json:"end,omitempty"`
	EndDate   time.Time `json:"-"`
}

type PaginatedRequest struct {
	PartnerID  string         `json:"partnerId"`
	MerchantID string         `json:"merchantID"`
	TerminalID string         `json:"terminalId,omitempty"`
	Periods    PeriodsRequest `json:"periods,omitempty"`
	Page       int64          `json:"page,omitempty"`
	Size       int64          `json:"size,omitempty"`
	Limit      int64          `json:"limit,omitempty"`
}
