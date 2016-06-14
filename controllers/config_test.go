package controllers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DaveBlooman/api-common/caching"
	"github.com/DaveBlooman/api-common/storage"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type ConfigControllerTestSuite struct {
	suite.Suite
	Cache      *caching.MockCache
	Storage    *storage.MockStorage
	Controller *ConfigController
	Router     *mux.Router
}

func (suite *ConfigControllerTestSuite) SetupTest() {
	suite.Cache = new(caching.MockCache)
	suite.Storage = new(storage.MockStorage)

	suite.Controller = &ConfigController{
		Cache:   suite.Cache,
		Storage: suite.Storage,
	}

	suite.Router = mux.NewRouter()
	suite.Router.Path("/v1/env/{env}/product/{product}/page/{page}").Methods("GET").HandlerFunc(suite.Controller.HandleGetRequest)
	suite.Router.Path("/v1/env/{env}/product/{product}/page/{page}").Methods("POST").HandlerFunc(suite.Controller.HandlePostRequest)
}

func (suite *ConfigControllerTestSuite) TearDownTest() {
	suite.Cache.AssertExpectations(suite.T())
	suite.Storage.AssertExpectations(suite.T())
}

func (suite *ConfigControllerTestSuite) TestGetFromCache() {
	url := "/v1/env/int/product/news/page/test"
	config := `{ "test": "config" }`

	suite.Cache.On("Get", url+".json").Return(config, nil).Once()

	r, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, r)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal("Redis Cache", w.HeaderMap["Content-Source"][0])
	suite.Equal(config, w.Body.String())
}

func (suite *ConfigControllerTestSuite) TestGetFromStorage() {
	url := "/v1/env/int/product/news/page/test"
	config := configExample

	suite.Cache.On("Get", url+".json").Return("", errors.New("error")).Once()
	var nilError *storage.Error
	suite.Storage.On("Get", url+".json").Return(config, nilError)
	suite.Cache.On("Set", url+".json", config, 60*time.Second).Return("", nil)

	r, _ := http.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, r)

	suite.Equal(http.StatusOK, w.Code)
	suite.Equal("AWS S3", w.HeaderMap["Content-Source"][0])
	suite.Equal(config, w.Body.String())
}

func (suite *ConfigControllerTestSuite) TestPost() {
	url := "/v1/env/int/product/news/page/test"
	config := configExample

	var nilError *storage.Error
	suite.Storage.On("Set", url+".json", config).Return(nilError)
	suite.Cache.On("Set", url+".json", config, 60*time.Second).Return("", nil)

	r, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(config)))
	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, r)

	suite.Equal(http.StatusOK, w.Code)
}

func TestConfigControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigControllerTestSuite))
}

var configExample = `
{
  "meta": {
    "title": "Config title",
    "description": "A longer description of the page layout configuration",
    "keywords": "config, keywords"
  },
  "components": [
    {
      "id": "component1",
      "endpoint": "http://hostname/path/to/component1",
      "must_succeed": true
    },
    {
      "id": "component2",
      "endpoint": "http://hostname/path/to/component2",
      "must_succeed": false
    }
  ]
}
`
