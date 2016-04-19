package data

import (
	"encoding/json"
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
	Temperature float32
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
	// Build a dictionary of data.
	// var data map[string]TimedTemperature

	var devices []Device
	db.Find(&devices)
	for index, device := range devices {
		log.Printf("Device #%d: %v", index, device)
	}

	// Select all temperature devices.
	/*
	   devices = db.select("SELECT DeviceID, Label FROM Devices WHERE Type IN ('ds18b20', 'dht22', 'dht11', 'am2302')")

	   # Iterate through the devices.
	   for device in devices:

	       # Get the latest temperature.
	       args = (device[0],)
	       rows = db.select("SELECT Timestamp, " + settings['t_unit'] + " FROM Temperature WHERE DeviceID = ? ORDER BY Timestamp DESC LIMIT 1", args)

	       # Fill in the data.
	       for row in rows:
	           timestamp = int(str(row[0]) + '000')
	           label = device[1]
	           temperature = row[1]
	           data[label] = [timestamp, temperature]

	   # Return as a string.
	   logger.info("/data/latest/temperature: {}".format(json.dumps(data)))
	   # /data/latest/temperature: {"28-000003ea01f5": [1461046802000, 4.38]}
	   return json.dumps(data)
	*/
	return `{"28-000003ea01f5": [1461046802000, 4.38]}`, nil
}

func GetLatestHumidity(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	return `{"28-000003ea01f5": [1461046802000, 4.38]}`, nil
}

func GetTemperatureDevices(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	var devices []Device
	data := make(map[string]string)
	db.Find(&devices)
	for _, device := range devices {
		data[fmt.Sprintf("%d", device.DeviceID)] = device.Label
	}
	j, err := json.Marshal(data)
	log.Printf("Devices: %s", j)
	var retErr util.HandlerError
	if err != nil {
		retErr = util.HandlerError{err, "Marshalling failed.", 500}
	}
	return j, &retErr
}

func GetHumidityDevices(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	return "HUMI", nil
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
