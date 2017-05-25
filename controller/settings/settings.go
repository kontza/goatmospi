package settings

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os/user"
	"strings"

	"github.com/gorilla/mux"
	rh "github.com/kontza/goatmospi/route_handler"
	"github.com/kontza/goatmospi/util"
)

// book model
type Settings struct {
	// How far into the past should data be loaded (in seconds)? Default to 1 week.
	RangeSeconds int `json:"range_seconds"`

	// The number of digits after the decimal place that will be stored.
	Precision int `json:"precision"`

	// Temperature unit of measure (C or F).
	TemperatureUnit string `json:"t_unit"`
}

var settings Settings

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

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
	// Expand environment: $HOME and '~' replaced by the actual path.
	homeDir := getHomeDir()
	settings.DB = strings.Replace(settings.DB, "$HOME", homeDir, -1)
	settings.DB = strings.Replace(settings.DB, "~", homeDir, -1)
	log.Printf("DB: %+v\n", settings.DB)
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
