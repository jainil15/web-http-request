package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	StatusCode int         `json:"status_code"`
	Result     interface{} `json:"result"`
	Message    string      `json:"message"`
}

type Error struct {
	StatusCode int         `json:"status_code"`
	Error      interface{} `json:"error"`
	Message    string      `json:"message"`
}

func ResponseHandler(w http.ResponseWriter, r *Response) {
	if r.StatusCode == 0 {
		r.StatusCode = http.StatusOK
	}
	jsonR, err := json.Marshal(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error encoding the response! -> %s\n", err)))
		return
	}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(r.StatusCode)
	w.Write(jsonR)
	return
}

func ErrorHandler(w http.ResponseWriter, e *Error) {
	if e.StatusCode == 0 {
		e.StatusCode = http.StatusInternalServerError
	}
	jsonR, err := json.Marshal(e)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error encoding the error! -> %s\n", err)))
		return
	}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(e.StatusCode)
	w.Write(jsonR)
	return
}
