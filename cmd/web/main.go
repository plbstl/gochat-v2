package main

import (
	"log"
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/paulebose/gochat-v2/internal/handlers"
)

var mux *pat.PatternServeMux

func init() {
	mux = pat.New()
	mux.Get("/", http.HandlerFunc(handlers.Home))

	// Register this pat with the default serve mux
	// so that other packages may also be exported.
	// (i.e. /debug/pprof/*)
	http.Handle("/", mux)
}

func main() {
	log.Println("Go Chat v2 \t http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalln("ListenAndServe:", err)
	}
}
