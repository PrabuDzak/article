package restapi

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (a *API) response(w http.ResponseWriter, statusCode int, response response) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (a *API) responseMessage(w http.ResponseWriter, statusCode int, message string) {
	a.response(w, statusCode, response{Message: message})
}

func (a *API) responseError(w http.ResponseWriter, statusCode int, err error) {
	a.response(w, statusCode, response{Message: err.Error()})
}
