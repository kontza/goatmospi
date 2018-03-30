package app_config

import (
	"path/filepath"
	"github.com/kontza/goatmospi/logger_factory"
	"github.com/jinzhu/configor"
)

// ApplicationConfig contains the application configuration.
type ApplicationConfig struct {
	Address    string `default:"127.0.0.1"`
	Port       string `default:"8080"`
	DirToServe string `default:"web"`
	Database struct {
		Hostname string `default:"localhost"`
		Port     string `default:"5432"`
		Name     string `default:"atmospi"`
		User     string `default:"atmospi"`
		Password string `default:"atmospi"`
	}
	Client struct {
		RangeSeconds    string `default:"604800"`
		Precision       string `default:"2"`
		TemperatureUnit string `default:"C"`
	}
}

var (
	appConfig = ApplicationConfig{}
	logger    = logger_factory.GetLogger()
)

/**
Read application config from the YAML.
*/
func ReadConfig(configFile string) ApplicationConfig {
	filename, _ := filepath.Abs(configFile)
	configor.Load(&appConfig, filename)
	printValue := appConfig
	printValue.Database.Password = "********"
	logger.Infof("Read config: %+v", printValue)
	return appConfig
}

/**
Return the Application Config -object.
 */
func GetApplicationConfig() ApplicationConfig {
	return appConfig
}
