package handlers

import (
	"encoding/json"
	"net/http"
)

type responseError struct {
	Code   int    `json:"code"`
	Error  string `json:"message,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type resp struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// response, err := json.Marshal(payload)
	response, err := json.Marshal(resp{code, "status ok", payload})
	if err != nil {
		response, _ = json.Marshal(responseError{Error: err.Error()})
		// code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func RespondWithJSONIndent(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.MarshalIndent(resp{code, "status ok", payload}, "", "  ")
	// response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		response, _ = json.Marshal(responseError{Error: err.Error()})
		code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {

	response, err := json.Marshal(resp{10000 + code, message, ""})
	if err != nil {
		response, _ = json.Marshal(responseError{Error: err.Error()})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)

	// RespondWithJSON(w, code, responseError{Code: code, Error: message})
}

func RespondWithDetailedError(w http.ResponseWriter, code int, message, detail string) {

	response, err := json.Marshal(resp{code, message + ":" + detail, ""})
	if err != nil {
		response, _ = json.Marshal(responseError{Error: err.Error()})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)

	// RespondWithJSON(w, code, responseError{Error: message, Detail: detail})
}

func RespondWithCode(w http.ResponseWriter, code int) {
	RespondWithJSON(w, code, "")

}
