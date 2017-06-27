package libs

import (
	"encoding/json"
	"net/http"
)

// Response is the result to user and it will be convert to a json object
type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

// createResponse function convert object to json response
func createResponse(result Response) []byte {
	res, _ := json.Marshal(result)
	return res
}

// Success function return a successful message to user
func Success(result interface{}) []byte {
	return createResponse(Response{
		StatusCode: 200,
		Message:    "ok",
		Data:       result,
	})
}

// Failed function return a failed message to user
func Failed(message string) []byte {
	return createResponse(Response{
		StatusCode: 500,
		Message:    message,
	})
}

// SendJSONResponseHeader 发送一个JSON的响应
func SendJSONResponseHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// Redirect 用于302跳转
func Redirect(url string, w http.ResponseWriter) {
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}

// SendNotFoundResponse 发送404响应
func SendNotFoundResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not Found!"))
}
