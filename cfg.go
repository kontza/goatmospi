package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	flags "github.com/jessevdk/go-flags"
	"github.com/juju/loggo"
	"gopkg.in/yaml.v2"
)

// App config.
type ApplicationConfig struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
	Client  struct {
		RangeSeconds    string `yaml:"range_seconds"`
		Precision       string `yaml:"precision"`
		TemperatureUnit string `yaml:"t_unit"`
	}
	Database struct {
		Hostname string `yaml:"hostname"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"db"`
}

// command line flags
var (
	options struct {
		Version bool   `short:"V" long:"version" description:"Show version information"`
		Verbose bool   `short:"v" long:"verbose" description:"Show verbose debug information"`
		File    string `short:"f" long:"file" default:"goatmospi.yml" description:"The configuration file (YAML) to use for the operations." value-name:"FILE"`
	}
	logger    = loggo.GetLogger("goatmospi")
	appConfig ApplicationConfig
)

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
	logger.Infof("Read YAML-config: %+v", printValue)
	return _config
}

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
		readConfig()
	}
}
