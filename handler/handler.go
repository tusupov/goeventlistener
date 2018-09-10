package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tusupov/goeventlistener/db"
)

// Error codes returned by failures
var (
	errNotFound = errors.New("404 page not found")
)

type handler struct {
	storage *db.Storage
}

func New(storage *db.Storage) *handler {
	return &handler{
		storage: storage,
	}
}

// NewListener handle function
func (h *handler) NewListener(w http.ResponseWriter, r *http.Request) {

	// Read body
	var listenerRequest db.ListenerRequest
	err := json.NewDecoder(r.Body).Decode(&listenerRequest)
	if err != nil {
		JSONErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Add listener
	err = h.storage.Add(listenerRequest)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Success result
	JSONSuccessResponse(w)

}

// DeleteListener handle function
func (h *handler) DeleteListener(w http.ResponseWriter, r *http.Request) {

	// Get vars
	vars := mux.Vars(r)
	listenerName := vars["listener"]

	// Delete listener by name
	err := h.storage.DeleteListener(listenerName)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Success result
	JSONSuccessResponse(w)

}

// PublishEvent handle function
func (h *handler) PublishEvent(w http.ResponseWriter, r *http.Request) {

	// Get vars
	vars := mux.Vars(r)
	eventName := vars["event"]

	// Event listeners
	// do http request
	err := h.storage.Publish(r.Context(), eventName)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Success result
	JSONSuccessResponse(w)

}

// NotFound handle function
func (h *handler) NotFound(w http.ResponseWriter, r *http.Request) {
	JSONErrorResponse(w, http.StatusNotFound, errNotFound.Error())
}
