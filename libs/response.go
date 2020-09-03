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
}

// Redirect 用于302跳转
func Redirect(url string, w http.ResponseWriter) {
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte("rediect to " + url))
}

// SendNotFoundResponse 发送404响应
func SendNotFoundResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("Not Found!"))
}

// SendFormInvalidResponse 发送表单不合法响应
func SendFormInvalidResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	_, _ = w.Write(createResponse(Response{
		StatusCode: 422,
		Message:    message,
	}))
}

// SendInternalServerErrorResponse 发送内部服务错误
func SendInternalServerErrorResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write(createResponse(Response{
		StatusCode: 500,
		Message:    message,
	}))
}
