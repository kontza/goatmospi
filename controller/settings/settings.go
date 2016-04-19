package settings

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	rh "github.com/kontza/goatmospi/route_handler"
	"github.com/kontza/goatmospi/util"
)

// book model
type Settings struct {
	// Absolute path to the SQLite database file.
	DB string `json:"db"`

	// How far into the past should data be loaded (in seconds)? Default to 1 week.
	RangeSeconds int `json:"range_seconds"`

	// The number of digits after the decimal place that will be stored.
	Precision int `json:"precision"`

	// Temperature unit of measure (C or F).
	TemperatureUnit string `json:"t_unit"`
}

var settings Settings

/**
 * Read settings from a JSON-file.
 * @param {[type]} yamlFile string [description]
 */
func LoadSettings(filename string) {
	file, e := ioutil.ReadFile(filename)
	if e != nil {
		log.Printf("File error: %v\n", e)
		return
	}
	json.Unmarshal(file, &settings)
	log.Printf("%+v\n", settings)
}

func GetSettingsData() Settings {
	return settings
}

func GetSettings(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	return settings, nil
}

func RegisterRoutes(router *mux.Router) {
	router.Handle("/settings", rh.RouteHandler(GetSettings)).Methods("GET")
}
