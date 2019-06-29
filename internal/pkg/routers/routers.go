package routers

import (
	"github.com/gorilla/mux"
	"net/http"
	controller "vegeta-kubernetes/internal/pkg/controllers"
)

type route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

// Routes is a list/slice of all routes for this micro-service
type Routes []route

// NewRouter return a router including method, path, name and handler
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandleFunc
		//handler = HandoeRouteWithLog(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	route{
		"GetMetrics",
		"GET",
		"/getMetrics",
		controller.GetMetrics,
	},
}
