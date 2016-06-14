package routes

import (
	"net/http"

	"github.com/DaveBlooman/api-common/caching"
	"github.com/DaveBlooman/api-common/storage"
	"github.com/DaveBlooman/bootstrap-api/controllers"
)

// Route options
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes slice
type Routes []Route

var statusController = &controllers.StatusController{}
var configController = &controllers.ConfigController{
	Cache:   &caching.RedisCache{},
	Storage: &storage.S3Storage{},
}

var routes = Routes{
	Route{
		"Status",
		"GET",
		"/status",
		statusController.HandleStatusRequest,
	},
	Route{
		"Product",
		"GET",
		"/v1/env/{env}/product/{product}/page/{page}",
		configController.HandleGetRequest,
	},
	Route{
		"Product",
		"POST",
		"/v1/env/{env}/product/{product}/page/{page}",
		configController.HandlePostRequest,
	},
}
