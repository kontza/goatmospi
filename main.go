package main

import (
	"fmt"
	htpl "html/template"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kontza/goatmospi/controller/data"
	"github.com/kontza/goatmospi/app_config"
	"github.com/kontza/goatmospi/controller/settings"
	flags "github.com/jessevdk/go-flags"
	"github.com/juju/loggo"
	"github.com/kontza/goatmospi/logger_factory"
)

var (
	// Command line flags.
	Options struct {
		Version bool   `short:"V" long:"version" description:"Show version information"`
		Verbose []bool   `short:"v" long:"verbose" description:"Show verbose debug information"`
		File    string `short:"f" long:"file" default:"goatmospi.yml" description:"The configuration file (YAML) to use for the operations." value-name:"FILE"`
	}
	logger= logger_factory.GetLogger()
	appConfig = app_config.ApplicationConfig{
		"127.0.0.1",
		"web",
		"4002",
		app_config.Client{
			"604800",
			"2",
			"C",
		},
		app_config.Database{
			"localhost",
			"dbname",
			"dbuser",
			"password",
		},
	}
)

// Build the main index.html.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	logger.Infof("Request URL: %v\n", r.URL)
	type TemplateContext struct {
		SubDomain       string
		OldestTimestamp int
		LatestTimestamp int
	}
	t, err := htpl.ParseFiles("web/index.html")
	if err != nil {
		logger.Infof("template.ParseFiles failed: %q", err)
	} else {
		prefixPath := r.Header.Get("Atmospi-Prefix-Path")
		oldest, newest := data.GetTimestampRange()
		err = t.Execute(w, &TemplateContext{prefixPath, oldest, newest})
		if err != nil {
			logger.Infof("t.Execute failed: %q", err)
		}
	}
}

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
	if logLevel<loggo.TRACE {
		logLevel = loggo.TRACE
	}
	logger.SetLogLevel(logLevel)
	logger.Debugf("Current level: %s", logger.LogLevel().String())
	if Options.Version {
		println("Goatmospi v1.0")
		return
	} else {
		logger.Infof("Initial YAML-config: %+v", appConfig)
		appConfig = app_config.ReadConfig(Options.File, appConfig)
		logger.Infof("Current YAML-config: %+v", appConfig)
	}

	// handle all requests by serving a file of the same name
	fs := http.Dir(appConfig.DirToServe)
	fileHandler := http.FileServer(fs)

	// setup routes
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler).Methods("GET")
	settings.RegisterRoutes(router)
	data.RegisterRoutes(router)
	router.PathPrefix("/static").Handler(handlers.LoggingHandler(os.Stderr, fileHandler))
	http.Handle("/", router)

	addr := fmt.Sprintf("%s:%s", appConfig.Address, appConfig.Port)
	logger.Infof("Listening on %s...", addr)
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(addr, nil)
	logger.Infof(err.Error())
}
