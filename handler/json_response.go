package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSONSuccessResponse
func JSONSuccessResponse(w http.ResponseWriter) {
	JSONResponse(w, http.StatusOK, map[string]string{"success": "ok"})
}

// JSONErrorResponse
// Wraps error string to json data
func JSONErrorResponse(w http.ResponseWriter, errorCode int, err string) {
	type ApiError struct {
		Error string `json:"error"`
	}
	JSONResponse(w, errorCode, ApiError{
		Error: err,
	})
}

// JSONResponse
// Write JSON data to http writer
// Encode resp to json and write to w
func JSONResponse(w http.ResponseWriter, statusCode int, resp interface{}) {

	respJson, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(respJson); err != nil {
		log.Println(err)
	}

}
