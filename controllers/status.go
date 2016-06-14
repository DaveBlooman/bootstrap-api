package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/DaveBlooman/api-common/logger"
)

// StatusController status endpoint controller
type StatusController struct{}

// HandleStatusRequest GET /status controller method
func (c *StatusController) HandleStatusRequest(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(&statusResponse{Status: "ok"})

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(response))

	logMessage := map[string]interface{}{"event": "get", "URI": "/status"}
	logger.Info(logMessage)
}

type statusResponse struct {
	Status string `json:"status" description:"status message"`
}
