package app_config

import (
	"path/filepath"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/kontza/goatmospi/logger_factory"
)

// ApplicationConfig contains the application configuration.
type ApplicationConfig struct {
	Address string `yaml:"address,omitempty"`
	DirToServe string `yaml:"dir_to_serve,omitempty"`
	Port    string `yaml:"port,omitempty"`
	Client
	Database `yaml:"db"`
}

// Client is the inner struct for a client configuration.
type Client struct {
	RangeSeconds    string `yaml:"range_seconds,omitempty"`
	Precision       string `yaml:"precision,omitempty"`
	TemperatureUnit string `yaml:"t_unit,omitempty"`
}

// Database is the inner struct for DB configuration.
type Database struct {
	Hostname string `yaml:"hostname,omitempty"`
	Name     string `yaml:"name,omitempty"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
}

var (
	appConfig ApplicationConfig
	logger    = logger_factory.GetLogger()
)

/**
Read application config from the YAML.
*/
func ReadConfig(configFile string, defaultConfig ApplicationConfig) ApplicationConfig {
	filename, _ := filepath.Abs(configFile)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var _config ApplicationConfig
	err = yaml.Unmarshal(yamlFile, &_config)
	if err != nil {
		panic(err)
	}

	printValue := _config
	printValue.Database.Password = "********"
	logger.Debugf("Read config: %+v", printValue)
	/*
		Copy initial values where read values were empty.
	*/
	if _config.Address == "" {
		_config.Address = defaultConfig.Address
	}
	if _config.Port == "" {
		_config.Port = defaultConfig.Port
	}
	if _config.Client.RangeSeconds == "" {
		_config.Client.RangeSeconds = defaultConfig.Client.RangeSeconds
	}
	if _config.Client.Precision == "" {
		_config.Client.Precision = defaultConfig.Client.Precision
	}
	if _config.Client.TemperatureUnit == "" {
		_config.Client.TemperatureUnit = defaultConfig.Client.TemperatureUnit
	}
	if _config.Database.Hostname == "" {
		_config.Database.Hostname = defaultConfig.Database.Hostname
	}
	if _config.Database.Name == "" {
		_config.Database.Name = defaultConfig.Database.Name
	}
	if _config.Database.User == "" {
		_config.Database.User = defaultConfig.Database.User
	}
	if _config.Database.Password == "" {
		_config.Database.Password = defaultConfig.Database.Password
	}
	appConfig = _config
	return _config
}

/**
Return the Application Config -object.
 */
func GetApplicationConfig() ApplicationConfig {
	return appConfig
}
