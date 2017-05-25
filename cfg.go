package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	flags "github.com/jessevdk/go-flags"
	"github.com/juju/loggo"
	"gopkg.in/yaml.v2"
)

// ApplicationConfig contains the application configuration.
type ApplicationConfig struct {
	Address string `yaml:"address,omitempty"`
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

// command line flags
var (
	options struct {
		Version bool   `short:"V" long:"version" description:"Show version information"`
		Verbose bool   `short:"v" long:"verbose" description:"Show verbose debug information"`
		File    string `short:"f" long:"file" default:"goatmospi.yml" description:"The configuration file (YAML) to use for the operations." value-name:"FILE"`
	}
	logger    = loggo.GetLogger("goatmospi")
	appConfig = ApplicationConfig{
		"127.0.0.1", "4002",
		Client{
			"604800",
			"2",
			"C",
		},
		Database{
			"localhost",
			"dbname",
			"dbuser",
			"password",
		},
	}
)

/**
Read application config from the YAML.
*/
func readConfig() ApplicationConfig {
	filename, _ := filepath.Abs(options.File)
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
	logger.Infof("Read config: %+v", printValue)
	/*
		Sanitize values.
	*/
	if _config.Address == "" {
		_config.Address = appConfig.Address
	}
	if _config.Port == "" {
		_config.Port = appConfig.Port
	}
	if _config.Client.RangeSeconds == "" {
		_config.Client.RangeSeconds = appConfig.Client.RangeSeconds
	}
	if _config.Client.Precision == "" {
		_config.Client.Precision = appConfig.Client.Precision
	}
	if _config.Client.TemperatureUnit == "" {
		_config.Client.TemperatureUnit = appConfig.Client.TemperatureUnit
	}
	if _config.Database.Hostname == "" {
		_config.Database.Hostname = appConfig.Database.Hostname
	}
	if _config.Database.Name == "" {
		_config.Database.Name = appConfig.Database.Name
	}
	if _config.Database.User == "" {
		_config.Database.User = appConfig.Database.User
	}
	if _config.Database.Password == "" {
		_config.Database.Password = appConfig.Database.Password
	}
	return _config
}

/**
Main routine.
*/
func main() {
	if _, err := flags.Parse(&options); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			println("Exit 1")
			os.Exit(1)
		}
	}

	if options.Verbose {
		logger.SetLogLevel(loggo.INFO)
	}
	if options.Version {
		println("Goatmospi v1.0")
	} else {
		appConfig = readConfig()
		logger.Infof("Initial YAML-config: %+v", appConfig)
	}
}
