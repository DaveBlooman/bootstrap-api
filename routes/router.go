package routes

import (
	"net/http"

	"github.com/DaveBlooman/api-common/logger"
	"github.com/gorilla/mux"
)

// APIRouter setup
func APIRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger.Log(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
