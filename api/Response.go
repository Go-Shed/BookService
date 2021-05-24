package api

type ApiResponse struct {
	HTTPCode int         `json:"http_code"`
	Data     interface{} `json:"data,omitempty"`
	Message  string      `json:"msg"`
}

type ErrorResponse struct {
	HTTPCode int    `json:"http_code"`
	Message  string `json:"msg"`
}
