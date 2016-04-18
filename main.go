package main

import (
	"flag"
	"fmt"
	htpl "html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kontza/goatmospi/controller/settings"
)

// Build the main index.html.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Request URL: %v\n", r.URL)
	type TemplateContext struct {
		SubDomain string
	}
	t, err := htpl.ParseFiles("web/index.html")
	if err != nil {
		fmt.Printf("template.ParseFiles failed: %q", err)
	} else {
		prefixPath := r.Header.Get("Atmospi-Prefix-Path")
		err = t.Execute(w, &TemplateContext{prefixPath})
		if err != nil {
			fmt.Printf("t.Execute failed: %q", err)
		}
	}
}

func main() {
	// command line flags
	port := flag.Int("port", 4002, "port to serve on")
	dir := flag.String("directory", "web", "directory of web files")
	flag.Parse()

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)

	// setup routes
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler).Methods("GET")
	settings.RegisterRoutes(router)
	router.PathPrefix("/static").Handler(handlers.LoggingHandler(os.Stderr, fileHandler))
	http.Handle("/", router)

	log.Printf("Running on port %d\n", *port)
	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}
