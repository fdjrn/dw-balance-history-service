package payload

type HistoryRequest struct {
	UID string `json:"uniqueId"`
	// 5, 10, 20, 50
	Limit int64 `json:"limit,omitempty"`
	// YearMonth
	Period int `json:"period,omitempty"`
}
