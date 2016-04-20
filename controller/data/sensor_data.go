package data

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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

func (tt *TimedTemperature) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%d, %f]", 1000*tt.Timestamp, tt.Temperature)), nil
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
		data[device.Label] = []string{fmt.Sprintf("%d", latest.Timestamp), fmt.Sprintf("%.2f", latestTemperature)}
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

func GetTemperatureData(w http.ResponseWriter, r *http.Request, deviceId int64, rangeMin float64, rangeMax float64) (interface{}, *util.HandlerError) {
	var temps []Temperature
	if rangeMin > 0 && rangeMax > 0 {
		db.Where("DeviceID = ? and Timestamp between ? and ?", deviceId, rangeMin, rangeMax).Order("Timestamp asc").Find(&temps)
	} else {
		db.Where("DeviceID = ?", deviceId).Order("Timestamp asc").Find(&temps)
	}
	data := []TimedTemperature{}
	var currentTemperature float64
	tUnit := settings.GetSettingsData().TemperatureUnit
	for _, temp := range temps {
		if tUnit == "C" {
			currentTemperature = temp.C
		} else {
			currentTemperature = temp.F
		}
		data = append(data, TimedTemperature{temp.Timestamp, currentTemperature})
	}
	count := 20
	ellipsis := "…"
	if len(data) < count {
		count = len(data)
		ellipsis = ""
	}
	log.Printf("Data: %d %v%s", len(data), data[:count], ellipsis)
	return data, nil
}

func GetHumidityData(w http.ResponseWriter, r *http.Request, deviceId int64, rangeMin float64, rangeMax float64) (interface{}, *util.HandlerError) {
	log.Printf("Data: nil")
	return 0, nil
}

func GetDeviceData(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	loadDatabase()
	vars := mux.Vars(r)
	log.Printf("Device data request: %v, %v", vars["deviceId"], vars["sensorType"])
	var retVal interface{}
	var err util.HandlerError
	rangeMin := 0.0
	rangeMax := 0.0
	queryArg := r.URL.Query().Get("range_min")
	var e error
	if len(queryArg) > 0 {
		rangeMin, e = strconv.ParseFloat(queryArg, 64)
		if e != nil {
			err = util.HandlerError{e, "range_min parse failed", 500}
			return nil, &err
		}
	}
	queryArg = r.URL.Query().Get("range_max")
	if len(queryArg) > 0 {
		rangeMax, e = strconv.ParseFloat(queryArg, 64)
		if e != nil {
			err = util.HandlerError{e, "range_max parse failed", 500}
			return nil, &err
		}
	}
	var deviceId int64
	deviceId, e = strconv.ParseInt(vars["deviceId"], 10, 32)
	if e != nil {
		err = util.HandlerError{e, "deviceId parse failed", 500}
		return nil, &err
	}
	var er *util.HandlerError
	switch vars["sensorType"] {
	case "temperature":
		retVal, er = GetTemperatureData(w, r, deviceId, rangeMin/1000.0, rangeMax/1000.0)
	case "humidity":
		retVal, er = GetHumidityData(w, r, deviceId, rangeMin/1000.0, rangeMax/1000.0)
	}
	return retVal, er
}

func RegisterRoutes(router *mux.Router) {
	router.Handle("/data/latest/{item}", rh.RouteHandler(GetLatest)).Methods("GET")
	router.Handle("/data/devices/{item}", rh.RouteHandler(GetDevices)).Methods("GET")
	router.Handle("/data/device/{deviceId}/{sensorType}", rh.RouteHandler(GetDeviceData)).Methods("GET")
}