package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"mystman.com/animated-pancake/internal/data"
	"mystman.com/animated-pancake/internal/service"
)

// Error variables
var (
	NotFound              = "data not found"
	MissingID             = "missing Id"
	InvalidInput          = "invalid input format"
	OperationNotSupported = "this operation is not supported"
)

// Handler - struct for handlers
type Handler struct {
	svc *service.Service
}

// NewHandler - return a new Handler
func NewHandler(s *service.Service) *Handler {
	return &Handler{
		svc: s,
	}
}

// HandleNetwork - handler for /network
func (h Handler) HandleNetwork(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleNetwork is called")

	switch r.Method {
	case http.MethodPost:
		b, _ := io.ReadAll(r.Body)

		var dat map[string]interface{}
		if err := json.Unmarshal(b, &dat); err != nil {
			errorResponse(w, InvalidInput, http.StatusBadRequest)
			return
		}

		dta, _ := h.svc.PostData(data.NetworkType, data.Data{Payload: dat})

		jsonReply := data.AsJSONString(dta)

		fmt.Fprintf(w, jsonReply)

	default:
		errorResponse(w, OperationNotSupported, http.StatusBadRequest)
	}
}

// HandleRoot - handle for /
func (h Handler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleRoot is called")

	path := strings.TrimPrefix(r.URL.Path, "/v1/")

	switch r.Method {

	// GET
	case http.MethodGet:
		if len(path) > 0 {
			dta, err := h.svc.GetData(path)
			if err != nil {
				errorResponse(w, NotFound, http.StatusNotFound)
				return
			}
			jsonResponse(w, dta, http.StatusOK)

		} else {
			query := r.URL.Query()
			paramID := query.Get("id")
			paramType := query.Get("type")

			dta, err := h.svc.GetAllData(paramID, paramType)
			if err != nil {
				errorResponse(w, NotFound, http.StatusNotFound)
				return
			}
			jsonResponse(w, dta, http.StatusOK)
		}

	// DELETE
	case http.MethodDelete:
		if len(path) == 0 {
			errorResponse(w, MissingID, http.StatusBadRequest)
			return
		}

		err := h.svc.DeleteData(path)
		if err != nil {
			errorResponse(w, NotFound, http.StatusNotFound)
			return
		}
	}

	resp := make(map[string]string)
	resp["status"] = "ok"
	jsonResponse(w, resp, http.StatusOK)

}

// HandleReadiness - handle readines and liveness probes
func HandleReadiness(w http.ResponseWriter, r *http.Request) {
	log.Printf("HandleReadiness is called")
	fmt.Fprintf(w, "Service is alive at: %v", r.URL)
}

//======================================================================================================

// jsonResponse - prepares a response from content as JSON
func jsonResponse(w http.ResponseWriter, content interface{}, httpStatusCode int) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	// jsonResp, _ := json.Marshal(content)
	jsonResp, _ := json.MarshalIndent(content, "", "  ")

	w.Write(jsonResp)
}

// errorResponse - prepares an error message response as JSON
func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	resp := make(map[string]string)
	resp["status"] = fmt.Sprintf("%v %v", strconv.Itoa(httpStatusCode), http.StatusText(httpStatusCode))
	resp["message"] = message

	jsonResp, _ := json.Marshal(resp)

	w.Write(jsonResp)
}
