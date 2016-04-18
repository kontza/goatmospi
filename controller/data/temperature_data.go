package data

import (
	// "fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	rh "github.com/kontza/goatmospi/route_handler"
	"github.com/kontza/goatmospi/util"
)

func GetLatestTemperature(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	return "20.0", nil
}

func GetLatestHumidity(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	return "65", nil
}

func GetTemperatureDevices(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	return "TEMP", nil
}

func GetHumidityDevices(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	return "HUMI", nil
}

func GetLatest(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
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
