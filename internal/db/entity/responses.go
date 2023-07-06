package entity

type Responses struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

//
//type ResponsePayloadData struct {
//	Total  int         `json:"total"`
//	Result interface{} `json:"results"`
//}

type PaginationInfo struct {
	PerPage     int64 `json:"perPage,omitempty"`
	CurrentPage int64 `json:"currentPage,omitempty"`
	LastPage    int64 `json:"lastPage,omitempty"`
}

type PaginatedDetailResponse struct {
	Total      int64          `json:"total,omitempty"`
	Result     interface{}    `json:"results,omitempty"`
	Pagination PaginationInfo `json:"pagination,omitempty"`
}

type PaginatedResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Data    PaginatedDetailResponse `json:"data,omitempty"`
}
