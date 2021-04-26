package api

type ApiResponse struct {
	ResponseCode int         `json:"response_code,omitempty"`
	Error        string      `json:"error,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}
