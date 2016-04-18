package main

import (
	"encoding/json"
	"flag"
	"fmt"
	htpl "html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kontza/goatmospi/controller/settings"
	"github.com/kontza/goatmospi/util"
)

// book model
type book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Id     int    `json:"id"`
}

// list of all of the books
var books = make([]book, 0)

// a custom type that we can use for handling errors and formatting responses
type handler func(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError)

// attach the standard ServeHTTP method to our handler so the http library can call it
func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// here we could do some prep work before calling the handler if we wanted to

	// call the actual handler
	response, err := fn(w, r)

	// check for errors
	if err != nil {
		log.Printf("ERROR: %v\n", err.Error)
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Message), err.Code)
		return
	}
	if response == nil {
		log.Printf("ERROR: response from method is nil\n")
		http.Error(w, "Internal server error. Check the logs.", http.StatusInternalServerError)
		return
	}

	// turn the response into JSON
	bytes, e := json.Marshal(response)
	if e != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	// send the response and log
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
	log.Printf("%s %s %s %d", r.RemoteAddr, r.Method, r.URL, 200)
}

var id = 0

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

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
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
	router.Handle("/settings", handler(settings.GetSettings)).Methods("GET")
	router.PathPrefix("/static").Handler(handlers.LoggingHandler(os.Stderr, fileHandler))
	http.Handle("/", router)

	log.Printf("Running on port %d\n", *port)
	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}
