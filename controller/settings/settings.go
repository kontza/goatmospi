package settings

import (
	rh "github.com/kontza/goatmospi/route_handler"
	"github.com/kontza/goatmospi/app_config"
	"github.com/gin-gonic/gin"
	"net/http"
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

func GetSettings(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, GetSettingsData())
}

func init() {
	rh.AddRouteHandler(rh.RoutingData{"GET", "/settings", GetSettings})
}
