package data

import (
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kontza/goatmospi/controller/settings"
	rh "github.com/kontza/goatmospi/route_handler"
	"github.com/kontza/goatmospi/logger_factory"
	"github.com/kontza/goatmospi/app_config"
	"math"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	db     *gorm.DB = nil
	logger          = logger_factory.GetLogger()
)

// TimedTemperature is a struct for timestamp and temperature pairs.
type TimedTemperature struct {
	Timestamp   int `json:"timestamp"`
	Temperature float64 `json:"temperature"`
}

func respondWithError(code int, message string, ctx *gin.Context) {
	resp := map[string]string{"error": message}
	ctx.JSON(code, resp)
	ctx.Abort()
}

// FromTemperatureC converts a Celsius Temperature into a TimedTemperature.
func FromTemperature(temperature Temperature) *TimedTemperature {
	retVal := new(TimedTemperature)
	retVal.Timestamp = 1000 * temperature.Timestamp
	if settings.GetSettingsData().TemperatureUnit == "C" {
		retVal.Temperature = temperature.C
	} else {
		retVal.Temperature = temperature.F
	}
	return retVal
}

// MarshalJSON converts a TimedTemperature into a JSON.
func (tt *TimedTemperature) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("[%d, %f]", 1000*tt.Timestamp, tt.Temperature)), nil
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}

func RoundPlus(f float64, places int) (float64) {
	shift := math.Pow(10, float64(places))
	return Round(f*shift) / shift;
}

func loadDatabase() {
	if db != nil {
		return
	}
	appConfig := app_config.GetApplicationConfig()
	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		appConfig.Database.Hostname,
		appConfig.Database.Port,
		appConfig.Database.Name,
		appConfig.Database.User,
		appConfig.Database.Password)
	var err error
	db, err = gorm.Open("postgres", connString)
	if err != nil {
		logger.Errorf("Connection to database failed: %s", err)
		panic("Cannot continue.")
	}
}

// GetTimestampRange gets the oldest and most recent timestamp from the data set in the database.
func GetTimestampRange() (int, int) {
	loadDatabase()
	var latest Temperature
	var oldest Temperature
	db.Order("Timestamp desc").Limit(1).Find(&latest)
	db.Order("Timestamp asc").Limit(1).Find(&oldest)
	return 1000 * oldest.Timestamp, 1000 * latest.Timestamp
}

// GetLatestTemperature gets the latest temperature from all the devices.
func GetLatestTemperature(ctx *gin.Context) {
	var devices []Device
	db.Where("Type in (?)", []string{"ds18b20", "dht22", "dht11", "am2302"}).Find(&devices)
	temperatureData := make(map[string][]float64)
	for _, device := range devices {
		var latest Temperature
		db.Where("DeviceID = ?", device.DeviceID).Order("Timestamp desc").Limit(1).Find(&latest)
		timedTemperature := FromTemperature(latest)
		temperatureData[device.Label] = []float64{float64(timedTemperature.Timestamp), RoundPlus(timedTemperature.Temperature, 2)}
	}
	ctx.JSON(http.StatusOK, temperatureData)
}

// GetLatestHumidity gets the latest humidity data from all the sensors.
func GetLatestHumidity(ctx *gin.Context) {
	var devices []Device
	db.Where("Type in (?)", []string{"dht22", "dht11", "am2302"}).Find(&devices)
	devicesData := make(map[string][]float64)
	for _, device := range devices {
		var latest Humidity
		db.Where("DeviceID = ?", device.DeviceID).Order("Timestamp desc").Limit(1).Find(&latest)
		devicesData[device.Label] = []float64{float64(1000 * latest.Timestamp), latest.H}
	}
	ctx.JSON(http.StatusOK, devicesData)
}

// GetTemperatureDevices returns all temperature devices in the database.
func GetTemperatureDevices(ctx *gin.Context) {
	var devices []Device
	devicesData := make(map[string]string)
	db.Where("Type in (?)", []string{"ds18b20", "dht22", "dht11", "am2302"}).Find(&devices)
	for _, device := range devices {
		devicesData[fmt.Sprintf("%d", device.DeviceID)] = device.Label
	}
	ctx.JSON(http.StatusOK, devicesData)
}

// GetHumidityDevices returns all humidity devices in the database.
func GetHumidityDevices(ctx *gin.Context) {
	var devices []Device
	devicesData := make(map[string]string)
	db.Where("Type in (?)", []string{"dht22", "dht11", "am2302"}).Find(&devices)
	for _, device := range devices {
		devicesData[fmt.Sprintf("%d", device.DeviceID)] = device.Label
	}
	ctx.JSON(http.StatusOK, devicesData)
}

// GetTemperatureData returns temperature data for the given time range.
func GetTemperatureData(ctx *gin.Context, deviceID int64, rangeMin float64, rangeMax float64) {
	var temps []Temperature
	if rangeMin > 0 && rangeMax > 0 {
		db.Where("DeviceID = ? and Timestamp between ? and ?", deviceID, rangeMin, rangeMax).Order("Timestamp asc").Find(&temps)
	} else {
		db.Where("DeviceID = ?", deviceID).Order("Timestamp asc").Find(&temps)
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
		data = append(data, TimedTemperature{temp.Timestamp, RoundPlus(currentTemperature, 2)})
	}
	count := 20
	ellipsis := "…"
	if len(data) < count {
		count = len(data)
		ellipsis = ""
	}
	logger.Infof("Data: %d %v%s", len(data), data[:count], ellipsis)
	ctx.JSON(http.StatusOK, data)
}

// GetHumidityData returns nil, since humidity sensors are not (yet) supported.
func GetHumidityData(ctx *gin.Context, deviceID int64, rangeMin float64, rangeMax float64) {
	logger.Infof("Data: nil")
	result := map[string]string{"status": "ok"}
	ctx.JSON(http.StatusOK, result)
}

// GetDeviceData returns a data set for the given device ID and sensor type.
func GetDeviceData(ctx *gin.Context) {
	loadDatabase()
	rangeMin := 0.0
	rangeMax := 0.0
	queryArg := ctx.Query("range_min")
	var e error
	if len(queryArg) > 0 {
		rangeMin, e = strconv.ParseFloat(queryArg, 64)
		if e != nil {
			respondWithError(500, "range_min parse failed", ctx)
			return
		}
	}
	queryArg = ctx.Query("range_max")
	if len(queryArg) > 0 {
		rangeMax, e = strconv.ParseFloat(queryArg, 64)
		if e != nil {
			respondWithError(500, "range_max parse failed", ctx)
		}
	}
	var deviceId int64
	deviceId, e = strconv.ParseInt(ctx.Param("deviceId"), 10, 32)
	if e != nil {
		respondWithError(http.StatusBadRequest, "deviceId parse filed", ctx)
	}
	switch ctx.Param("sensorType") {
	case "temperature":
		GetTemperatureData(ctx, deviceId, rangeMin/1000.0, rangeMax/1000.0)
	case "humidity":
		GetHumidityData(ctx, deviceId, rangeMin/1000.0, rangeMax/1000.0)
	default:
		respondWithError(http.StatusBadRequest, "invalid sensorType", ctx)
	}
}

// Register routes that this module takes care of.
func init() {
	rh.AddRouteHandler(rh.RoutingData{"GET", "/data/latest/temperature", GetLatestTemperature})
	rh.AddRouteHandler(rh.RoutingData{"GET", "/data/latest/humidity", GetLatestHumidity})
	rh.AddRouteHandler(rh.RoutingData{"GET", "/data/devices/temperature", GetTemperatureDevices})
	rh.AddRouteHandler(rh.RoutingData{"GET", "/data/devices/humidity", GetHumidityDevices})
	rh.AddRouteHandler(rh.RoutingData{"GET", "/data/device/:deviceId/:sensorType", GetDeviceData})
}
