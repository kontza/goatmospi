package main

import (
	"fmt"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	flags "github.com/jessevdk/go-flags"
	"github.com/juju/loggo"
	"github.com/kontza/goatmospi/app_config"
	"github.com/kontza/goatmospi/logger_factory"
	rh "github.com/kontza/goatmospi/route_handler"
	// Load the modules so that they can register their routes.
	_ "github.com/kontza/goatmospi/controller/data"
	_ "github.com/kontza/goatmospi/controller/index"
	_ "github.com/kontza/goatmospi/controller/settings"
)

var (
	// Command line flags.
	Options struct {
		Version bool   `short:"V" long:"version" description:"Show version information"`
		Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
		File    string `short:"c" long:"config" default:"goatmospi.yml" description:"The configuration file (YAML) to use for the operations." value-name:"FILE"`
	}
	logger    = logger_factory.GetLogger()
	appConfig = app_config.ApplicationConfig{}
)

func main() {
	if _, err := flags.Parse(&Options); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			println("Exit 1")
			os.Exit(1)
		}
	}
	// Increase verboseness.
	logLevel := loggo.WARNING
	for range Options.Verbose {
		logLevel--
	}
	if logLevel < loggo.TRACE {
		logLevel = loggo.TRACE
	}
	logger.SetLogLevel(logLevel)
	logger.Debugf("Current level: %s", logger.LogLevel().String())
	if Options.Version {
		println("Goatmospi v1.0")
		return
	} else {
		appConfig = app_config.ReadConfig(Options.File) // , appConfig)
	}

	// Set up Gin.
	ginRouter := gin.Default()
	ginRouter.StaticFile("/favicon.png", "./resources/favicon.png")
	ginRouter.Static("/static", path.Join(appConfig.DirToServe, "static"))
	ginRouter.LoadHTMLGlob("templates/*")

	// setup routes
	logger.Infof("Setting up routes...")
	rh.RegisterRoutes(ginRouter)

	// Start the server.
	addr := fmt.Sprintf("%s:%s", appConfig.Address, appConfig.Port)
	logger.Infof("Listening on %s...", addr)
	ginRouter.Run(addr)
}
