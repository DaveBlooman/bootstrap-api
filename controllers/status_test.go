package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandleStatusRequestSuccess(t *testing.T) {
	controller := &StatusController{}

	assert := assert.New(t)

	r, _ := http.NewRequest("GET", "/status", nil)

	router := mux.NewRouter()
	router.HandleFunc("/status", controller.HandleStatusRequest)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	var response interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	responseMap := response.(map[string]interface{})

	assert.Equal(responseMap["status"], "ok")
	assert.Equal(w.Code, http.StatusOK, "status code")
}
