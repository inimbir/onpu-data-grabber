package http

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type SearchRequest struct {
	Type  int    `json:"type"`
	Group string `json:"group"`
}

type SearchResonse struct {
	Message string `json:"message"`
}

const (
	FullSearch = 1
	TagSearch  = 2
)

func Run() {
	http.HandleFunc("/search", search)
	http.HandleFunc("/search-traditional", searchTraditional)
	http.HandleFunc("/search-criterias", searchHashTags)
	http.HandleFunc("/criterias", manageHashTags)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func GetLog() *log.Logger {
	// set location of log file
	var logpath = "output"
	flag.Parse()
	var file, err1 = os.Create(logpath)
	if err1 != nil {
		panic(err1)
	}
	Log := log.New(file, "", log.LstdFlags|log.Lshortfile)
	return Log
}
