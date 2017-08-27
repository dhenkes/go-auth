package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func resData(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	js, _ := json.Marshal(response)
	fmt.Fprint(w, string(js), "\n")
}

func resError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := Error{Code: code, Message: message}
	js, _ := json.Marshal(response)
	fmt.Fprint(w, string(js), "\n")
}

func resCouldNotReadBody(w http.ResponseWriter) {
	resError(w, http.StatusBadRequest, "Could not parse request body")
}

func resCouldNotParseBody(w http.ResponseWriter) {
	resError(w, http.StatusBadRequest, "Could not parse request body")
}

func resCouldNotGenerateToken(w http.ResponseWriter) {
	resError(w, http.StatusInternalServerError, "Internal server error")
}

func resCouldNotGenerateUUID(w http.ResponseWriter) {
	resError(w, http.StatusInternalServerError, "Internal server error")
}

func resCouldNotPrepareStmt(w http.ResponseWriter) {
	resError(w, http.StatusInternalServerError, "Internal server error")
}

func resCouldNotInsertIntoDB(w http.ResponseWriter) {
	resError(w, http.StatusInternalServerError, "Internal server error")
}

func resCouldNotHashPassword(w http.ResponseWriter, err error) {
	resError(w, http.StatusInternalServerError, "Internal server error")
}

func resNoRowFound(w http.ResponseWriter) {
	resError(w, http.StatusNotFound, "Not found")
}
