package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tusupov/goeventlistener/db"
	"github.com/tusupov/goeventlistener/handler"
	"github.com/tusupov/goeventlistener/middleware"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	port = kingpin.Flag("port", "Listening port").Short('p').Default("8080").Int()
)

func init() {
	kingpin.Parse()
}

func main() {

	localStorage := db.New()

	r := mux.NewRouter()

	// Handle function
	h := handler.New(localStorage)
	r.HandleFunc("/listener", h.NewListener).Methods(http.MethodPost)
	r.HandleFunc("/listener/{listener}", h.DeleteListener).Methods(http.MethodDelete)
	r.HandleFunc("/publish/{event}", h.PublishEvent).Methods(http.MethodPost)

	// 404
	r.NotFoundHandler = http.HandlerFunc(h.NotFound)

	// Middleware
	middleware.SetLogger(os.Stderr)
	r.Use(middleware.Panic, middleware.AccessLog)

	// Start server
	log.Printf("Listening port [%d] ...", *port)
	if err := http.ListenAndServe(":"+strconv.Itoa(*port), r); err != nil {
		log.Fatal(err)
	}

}
