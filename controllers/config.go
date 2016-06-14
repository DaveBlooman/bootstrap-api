package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/DaveBlooman/api-common/caching"
	"github.com/DaveBlooman/api-common/logger"
	"github.com/DaveBlooman/api-common/storage"
	"github.com/DaveBlooman/bootstrap-api/appconfig"
	"github.com/DaveBlooman/bootstrap-api/errors"
	"github.com/gorilla/mux"
)

var config appconfig.Config

func init() {
	config, _ = appconfig.LoadConfig()
}

// ConfigController config endpoint controller
type ConfigController struct {
	Cache   caching.Cache
	Storage storage.Storage
}

// HandleGetRequest GET /v1/env/{env}/product/{product}/page/{page} controller method
func (c *ConfigController) HandleGetRequest(w http.ResponseWriter, r *http.Request) {
	var logMessage map[string]interface{}

	vars := mux.Vars(r)
	env, product, page := vars["env"], vars["product"], vars["page"]
	url := fmt.Sprintf("/v1/env/%s/product/%s/page/%s.json", env, product, page)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var config string

	cacheResponse, cacheErr := c.Cache.Get(url)
	if cacheErr != nil {
		s3Response, s3Err := c.Storage.Get(url)
		config = s3Response

		if s3Err != nil {
			customerrors.Handle(w, r, &customerrors.Controller{
				Event:   "S3GetError",
				Status:  s3Err.Status,
				Message: s3Err.Message,
			})
			return
		}

		c.setInRedis(w, url, config)

		w.Header().Set("Content-Source", "AWS S3")

		logMessage = map[string]interface{}{"event": "S3GetSuccess", "uri": r.URL}
	} else {
		config = cacheResponse

		w.Header().Set("Content-Source", "Redis Cache")

		logMessage = map[string]interface{}{"event": "CacheGetSuccess", "uri": r.URL}
	}

	w.Write([]byte(config))

	logger.Info(logMessage)
}

// HandlePostRequest POST /v1/env/{env}/product/{product}/page/{page} controller method
func (c *ConfigController) HandlePostRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	env, product, page := vars["env"], vars["product"], vars["page"]
	url := fmt.Sprintf("/v1/env/%s/product/%s/page/%s.json", env, product, page)

	var logMessage map[string]interface{}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		customerrors.Handle(w, r, &customerrors.Controller{
			Event:   "InvalidPostRequest",
			Message: err.Error(),
		})
		return
	}

	putErr := c.Storage.Set(url, string(data))
	if putErr != nil {
		customerrors.Handle(w, r, &customerrors.Controller{
			Status:  500,
			Event:   "S3PutError",
			Message: putErr.Message,
		})
		return
	}

	c.setInRedis(w, url, string(data))

	response, _ := json.Marshal(&postConfigResponse{Message: "config added successfully"})
	w.Write([]byte(response))

	logMessage = map[string]interface{}{"event": "S3PutSuccess", "uri": url}
	logger.Info(logMessage)
}

func (c *ConfigController) setInRedis(w http.ResponseWriter, url string, data string) {
	_, err := c.Cache.Set(url, data, 60*time.Second)

	if err != nil {
		logMessage := map[string]interface{}{
			"event":   "RedisConnectionError",
			"message": err,
		}
		logger.Error(logMessage)
	}
}

type postConfigResponse struct {
	Message string `json:"message" description:"Success message"`
}
