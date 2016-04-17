package main

import (
	"encoding/json"
	"flag"
	"fmt"
	htpl "html/template"
	"io/ioutil"
	"log"
	"net/http"
	// "os"
	"strconv"

	// "github.com/gorilla/handlers"
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

func listBooks(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	return books, nil
}

func getBook(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	// mux.Vars grabs variables from the path
	param := mux.Vars(r)["id"]
	id, e := strconv.Atoi(param)
	if e != nil {
		return nil, &util.HandlerError{e, "Id should be an integer", http.StatusBadRequest}
	}
	b, index := getBookById(id)

	if index < 0 {
		return nil, &util.HandlerError{nil, "Could not find book " + param, http.StatusNotFound}
	}

	return b, nil
}

func parseBookRequest(r *http.Request) (book, *util.HandlerError) {
	// the book payload is in the request body
	data, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return book{}, &util.HandlerError{e, "Could not read request", http.StatusBadRequest}
	}

	// turn the request body (JSON) into a book object
	var payload book
	e = json.Unmarshal(data, &payload)
	if e != nil {
		return book{}, &util.HandlerError{e, "Could not parse JSON", http.StatusBadRequest}
	}

	return payload, nil
}

func addBook(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	payload, e := parseBookRequest(r)
	if e != nil {
		return nil, e
	}

	// it's our job to assign IDs, ignore what (if anything) the client sent
	payload.Id = getNextId()
	books = append(books, payload)

	// we return the book we just made so the client can see the ID if they want
	return payload, nil
}

func updateBook(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	payload, e := parseBookRequest(r)
	if e != nil {
		return nil, e
	}

	_, index := getBookById(payload.Id)
	books[index] = payload
	return make(map[string]string), nil
}

func removeBook(w http.ResponseWriter, r *http.Request) (interface{}, *util.HandlerError) {
	param := mux.Vars(r)["id"]
	id, e := strconv.Atoi(param)
	if e != nil {
		return nil, &util.HandlerError{e, "Id should be an integer", http.StatusBadRequest}
	}
	// this is jsut to check to see if the book exists
	_, index := getBookById(id)

	if index < 0 {
		return nil, &util.HandlerError{nil, "Could not find entry " + param, http.StatusNotFound}
	}

	// remove a book from the list
	books = append(books[:index], books[index+1:]...)
	return make(map[string]string), nil
}

// searches the books for the book with `id` and returns the book and it's index, or -1 for 404
func getBookById(id int) (book, int) {
	for i, b := range books {
		if b.Id == id {
			return b, i
		}
	}
	return book{}, -1
}

var id = 0

// increments id and returns the value
func getNextId() int {
	id += 1
	return id
}

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
	dir := flag.String("directory", "web/static", "directory of web files")
	flag.Parse()

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)

	// setup routes
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler).Methods("GET")
	router.Handle("/settings", handler(settings.GetSettings)).Methods("GET")
	router.Handle("/books", handler(listBooks)).Methods("GET")
	router.Handle("/books", handler(addBook)).Methods("POST")
	router.Handle("/books/{id}", handler(getBook)).Methods("GET")
	router.Handle("/books/{id}", handler(updateBook)).Methods("POST")
	router.Handle("/books/{id}", handler(removeBook)).Methods("DELETE")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", loggingHandler(fileHandler)))
	http.Handle("/", router)

	// bootstrap some data
	books = append(books, book{"Ender's Game", "Orson Scott Card", getNextId()})
	books = append(books, book{"Code Complete", "Steve McConnell", getNextId()})
	books = append(books, book{"World War Z", "Max Brooks", getNextId()})
	books = append(books, book{"Pragmatic Programmer", "David Thomas", getNextId()})

	log.Printf("Running on port %d\n", *port)

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}
