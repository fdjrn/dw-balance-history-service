package handlers

type ResponsePayload struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    ResponsePayloadData `json:"data,omitempty"`
}

type ResponsePayloadData struct {
	Total  int         `json:"total"`
	Result interface{} `json:"results"`
}

type ResponsePayloadDataPaginated struct {
	Result      interface{} `json:"results,omitempty"`
	Total       int64       `json:"total,omitempty"`
	PerPage     int64       `json:"perPage,omitempty"`
	CurrentPage int64       `json:"currentPage,omitempty"`
	LastPage    int64       `json:"lastPage,omitempty"`
}

type ResponsePayloadPaginated struct {
	Success bool                         `json:"success"`
	Message string                       `json:"message"`
	Data    ResponsePayloadDataPaginated `json:"data,omitempty"`
}
