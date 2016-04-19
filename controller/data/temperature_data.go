package data

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kontza/goatmospi/controller/settings"
	rh "github.com/kontza/goatmospi/route_handler"
	"github.com/kontza/goatmospi/util"
)

var db *gorm.DB

type TimedTemperature struct {
	Timestamp   int
	Temperature float64
}

func loadDatabase() {
	var err error
	if db == nil {
		db, err = gorm.Open("sqlite3", settings.GetSettingsData().DB)
		if err != nil {
			log.Fatalf("DB open failed: %v", err)
		} else {
			log.Printf("Opened the DB '%v'.", settings.GetSettingsData().DB)
		}
	} else {
		log.Printf("Reused the existing database connection.")
	}
}

func GetLatestTemperature(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	var devices []Device
	tUnit := settings.GetSettingsData().TemperatureUnit
	db.Where("Type in (?)", []string{"ds18b20", "dht22", "dht11", "am2302"}).Find(&devices)
	data := make(map[string][]string)
	var latestTemperature float64
	for _, device := range devices {
		var latest Temperature
		db.Where("DeviceID = ?", device.DeviceID).Order("Timestamp desc").Limit(1).Find(&latest)
		if tUnit == "C" {
			latestTemperature = latest.C
		} else {
			latestTemperature = latest.F
		}
		data[device.Label] = []string{fmt.Sprintf("%d", latest.Timestamp), fmt.Sprintf("%f", latestTemperature)}
	}
	return data, nil
}

func GetLatestHumidity(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	var devices []Device
	db.Where("Type in (?)", []string{"dht22", "dht11", "am2302"}).Find(&devices)
	data := make(map[string][]string)
	for _, device := range devices {
		var latest Humidity
		db.Where("DeviceID = ?", device.DeviceID).Order("Timestamp desc").Limit(1).Find(&latest)
		data[device.Label] = []string{fmt.Sprintf("%d", latest.Timestamp), fmt.Sprintf("%f", latest.H)}
	}
	return data, nil
}

func GetTemperatureDevices(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	var devices []Device
	data := make(map[string]string)
	db.Where("Type in (?)", []string{"ds18b20", "dht22", "dht11", "am2302"}).Find(&devices)
	for _, device := range devices {
		data[fmt.Sprintf("%d", device.DeviceID)] = device.Label
	}
	return data, nil
}

func GetHumidityDevices(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	var devices []Device
	data := make(map[string]string)
	db.Where("Type in (?)", []string{"dht22", "dht11", "am2302"}).Find(&devices)
	for _, device := range devices {
		data[fmt.Sprintf("%d", device.DeviceID)] = device.Label
	}
	return data, nil
}

func GetLatest(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	loadDatabase()
	vars := mux.Vars(r)
	item := vars["item"]
	log.Printf("Latest request: %v", item)
	var retVal interface{}
	var err *util.HandlerError
	switch item {
	case "temperature":
		retVal, err = GetLatestTemperature(w, r)
	case "humidity":
		retVal, err = GetLatestHumidity(w, r)
	}
	return retVal, err
}

func GetDevices(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	loadDatabase()
	vars := mux.Vars(r)
	item := vars["item"]
	log.Printf("Devices request: %v", item)
	var retVal interface{}
	var err *util.HandlerError
	switch item {
	case "temperature":
		retVal, err = GetTemperatureDevices(w, r)
	case "humidity":
		retVal, err = GetHumidityDevices(w, r)
	}
	return retVal, err
}

func RegisterRoutes(router *mux.Router) {
	router.Handle("/data/latest/{item}", rh.RouteHandler(GetLatest)).Methods("GET")
	router.Handle("/data/devices/{item}", rh.RouteHandler(GetDevices)).Methods("GET")
}
