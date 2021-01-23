package checkbodyplugin

import (
	"fmt"
	"net/http"
)

//ResponseError contains a failuer message
type ResponseError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
	Raw     string `json:"raw,omitempty"`
}

//GetMessage get response message
func (resErr *ResponseError) GetMessage() string {
	var s string
	if resErr.Raw == "" {
		s = fmt.Sprintf(`{
			"data": null,
			"error": {
				"code": "%s",
				"message": "%s"
			}
		}`, resErr.Code, resErr.Message)
	} else {
		s = resErr.Raw
	}
	return s
}

//Response write error response to client
func (resErr *ResponseError) Response(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(resErr.Status)
	rw.Write([]byte(resErr.GetMessage()))
}
