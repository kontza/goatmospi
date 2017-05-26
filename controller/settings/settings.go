package settings

import (
	"log"
	"net/http"
	"os/user"

	"github.com/gorilla/mux"
	rh "github.com/kontza/goatmospi/route_handler"
	"github.com/kontza/goatmospi/util"
	"github.com/kontza/goatmospi/app_config"
)

var settings app_config.Client

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func GetSettingsData() app_config.Client {
	return settings
}

func GetSettings(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	return settings, nil
}

func RegisterRoutes(router *mux.Router) {
	router.Handle("/settings", rh.RouteHandler(GetSettings)).Methods("GET")
}
