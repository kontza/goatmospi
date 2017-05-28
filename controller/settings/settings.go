package settings

import (
	"net/http"

	"github.com/gorilla/mux"
	rh "github.com/kontza/goatmospi/route_handler"
	"github.com/kontza/goatmospi/util"
	"github.com/kontza/goatmospi/app_config"
)

type ClientSettings struct {
	// How far into the past should data be loaded (in seconds)? Default to 1 week.
	RangeSeconds string `json:"range_seconds"`

	// The number of digits after the decimal place that will be stored.
	Precision string `json:"precision"`

	// Temperature unit of measure (C or F).
	TemperatureUnit string `json:"t_unit"`
}

func GetSettingsData() ClientSettings {
	currentConfig := app_config.GetApplicationConfig().Client
	return ClientSettings{currentConfig.RangeSeconds,
	currentConfig.Precision,
	currentConfig.TemperatureUnit}
}

func GetSettings(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	return GetSettingsData(), nil
}

func RegisterRoutes(router *mux.Router) {
	router.Handle("/settings", rh.RouteHandler(GetSettings)).Methods("GET")
}
