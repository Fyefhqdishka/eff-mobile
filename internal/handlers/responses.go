package handlers

import (
	"encoding/json"
	"net/http"
)

const (
	statusOK  = "OK"
	statusErr = "Error"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  any    `json:"result"`
}

func SendSuccess(result any) Response {
	return Response{
		Status: statusOK,
		Result: result,
	}
}

func SendError(msg string) Response {
	return Response{
		Status:  statusErr,
		Message: msg,
	}
}

func (h *Handlers) response(w http.ResponseWriter, r Response, statusCode int) {
	data, err := json.Marshal(r)
	if err != nil {
		msg := "can't marshal response"
		h.log.Error(msg, ", err=", err)
		r = SendError(msg)
		statusCode = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}
