package utils

import (
	"encoding/json"
	"net/http"
)

type JsonHTTPResponse struct {
	Status     string      `json:"status"`
	StatusCode int         `json:"statusCode"`
}

type JsonHTTPResponseData struct {
	Data       interface{} `json:"data"`
	JsonHTTPResponse
}

type JsonHTTPResponseError struct {
	Error string `json:"error"`
	JsonHTTPResponse
}

type JsonHTTPResponseMessage struct {
	Message string `json:"msg"`
	JsonHTTPResponse
}

func WriteJSONResponseData (w http.ResponseWriter, code int, data interface{}) {
	res := &JsonHTTPResponseData{
		Data: data,
		JsonHTTPResponse: JsonHTTPResponse{
			Status: http.StatusText(code),
			StatusCode: code,
		},
	}

	writeJSONResponse(w, code, res)
}

func WriteJSONResponseError (w http.ResponseWriter, code int, err error) {
	res := &JsonHTTPResponseError{
		Error: err.Error(),
		JsonHTTPResponse: JsonHTTPResponse{
			Status: http.StatusText(code),
			StatusCode: code,
		},
	}

	writeJSONResponse(w, code, res)
}

func WriteJSONResponseMessage (w http.ResponseWriter, code int, message string) {
	res := &JsonHTTPResponseMessage{
		Message: message,
		JsonHTTPResponse: JsonHTTPResponse{
			Status: http.StatusText(code),
			StatusCode: code,
		},
	}

	writeJSONResponse(w, code, res)
}

func writeJSONResponse(w http.ResponseWriter, code int, response interface{}) {
	w.Header().Set("Context-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

