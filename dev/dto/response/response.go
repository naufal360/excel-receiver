package response

type BaseResponse struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
	Detail    Detail `json:"detail"`
}

type Detail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
