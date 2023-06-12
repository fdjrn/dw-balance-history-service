package payload

type HistoryRequest struct {
	UID    string               `json:"uniqueId"`
	Limit  int64                `json:"limit,omitempty"` // 5, 10, 20, 50
	Period HistoryRequestPeriod `json:"period,omitempty"`
}

type HistoryRequestPeriod struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}
