package route_handler

import (
	"net/http"
	"github.com/kontza/goatmospi/util"
	"log"
	"github.com/gin-gonic/gin"
	"strings"
)

type RoutingData struct {
	Method  string
	Route   string
	Handler gin.HandlerFunc
}

// a custom type that we can use for handling errors and formatting responses
type RouteHandler func(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError)

var (
	routes []RoutingData
)

/**
AddRouteHandler adds a new handler to the array of registered route handlers.
 */
func AddRouteHandler(newEntry RoutingData) {
	log.Printf(">>> Route handler added for %s", newEntry.Route)
	routes = append(routes, newEntry)
}

/**
GetRouteHandlers returns the array of currently registered route handlers.
 */
func GetRouteHandlers() []RoutingData {
	return routes
}

/**
Add the registered route handlers to the given Gin router.
 */
func RegisterRoutes(r *gin.Engine) {
	for _, route := range routes {
		switch strings.ToLower(route.Method) {
		case "get":
			r.GET(route.Route, route.Handler)
			break
		case "post":
			r.POST(route.Route, route.Handler)
			break
		}
	}
}