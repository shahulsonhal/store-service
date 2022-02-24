package app

import (
	"encoding/json"
	"net/http"
)

const (
	statusOK   = "ok"
	statusFail = "nok"
)

// Response implements standard JSON response payload structure.
type Response struct {
	Status     string          `json:"status"`
	Result     json.RawMessage `json:"result,omitempty"`
	Error      *ResponseError  `json:"error,omitempty"`
	Pagination json.RawMessage `json:"pagination,omitempty"`
}

// ResponseError implements the standard Error response structure.
type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// send sends a successful JSON response using
// the standard success format
func send(w http.ResponseWriter, status int, result interface{}) {
	rj, err := json.Marshal(result)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	r := &Response{
		Status: statusOK,
		Result: rj,
	}

	j, err := json.Marshal(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)
}

// fail ends an unsuccessful JSON response with the standard failure format.
func fail(w http.ResponseWriter, status int, msg string) {
	r := &Response{
		Status: statusFail,
		Error: &ResponseError{
			Code:    status,
			Message: msg,
		},
	}
	j, err := json.Marshal(r)
	if err != nil {
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)
}
