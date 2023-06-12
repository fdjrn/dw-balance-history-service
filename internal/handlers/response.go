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
