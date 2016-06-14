package customerrors

import (
	"encoding/json"
	"net/http"

	"github.com/DaveBlooman/api-common/logger"
)

// Handle status
func Handle(w http.ResponseWriter, r *http.Request, e *Controller) {
	status := e.Status
	if status == 0 {
		status = 500
	}

	response, _ := json.Marshal(&Response{
		Error:   e.Event,
		Message: e.Message,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))

	logMessage := map[string]interface{}{
		"event":   e.Event,
		"message": e.Message,
		"uri":     r.URL.Path,
	}
	logger.Error(logMessage)
}

// Response items
type Response struct {
	Error   string      `json:"error" description:"Error identifier"`
	Message interface{} `json:"message" description:"Error description"`
}

// Controller options
type Controller struct {
	Status  int
	Event   string
	Message interface{}
}
