package utils

// response format
type HTTPResponse struct {
	Status     string      `json:"status"`
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
}

// response error format
type HTTPErrorMessage struct {
	Error string `json:"error"`
}

// response info format
type HTTPInfoMessage struct {
	Message string `json:"message"`
}
