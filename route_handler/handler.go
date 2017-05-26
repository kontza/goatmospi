package route_handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kontza/goatmospi/util"
	"github.com/kontza/goatmospi/logger_factory"
)

var (
	logger = logger_factory.GetLogger()
)

// a custom type that we can use for handling errors and formatting responses
type RouteHandler func(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError)

// attach the standard ServeHTTP method to our handler so the http library can call it
func (fn RouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// here we could do some prep work before calling the handler if we wanted to

	// call the actual handler
	response, err := fn(w, r)

	// check for errors
	if err != nil {
		logger.Infof("ERROR: %v\n", err.Error)
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Message), err.Code)
		return
	}
	if response == nil {
		logger.Infof("ERROR: response from method is nil.\n")
		http.Error(w, "Internal server error. Check the logs.", http.StatusInternalServerError)
		return
	}

	// turn the response into JSON
	bytes, e := json.Marshal(response)
	if e != nil {
		http.Error(w, "Error marshalling JSON.", http.StatusInternalServerError)
		return
	}

	// send the response and log
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
	count := 256
	ellipsis := "…"
	if len(bytes) < count {
		count = len(bytes)
		ellipsis = ""
	}
	logger.Infof("%s %s %s %d %s%s", r.RemoteAddr, r.Method, r.URL, 200, bytes[:count], ellipsis)
}
