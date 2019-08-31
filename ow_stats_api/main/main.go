package main

import (
	"fmt"
	"log"
	"net/http"

	"../accountstats"
	"github.com/gorilla/mux"
)

var (
	testbtag = []string{"brokenglass-4115,orisa"}
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, accountstats.ConcurrentGetRawAccountStats(testbtag)[0])
	fmt.Println("Endpoint Hit: homePage")
}

func getBasicStats(w http.ResponseWriter, r *http.Request) {
	vars = mux.Vars(r)

}

// HandleRequests maps API requests to the correct function
func HandleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/OWAccountStats/{btag}", getBasicStats)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	fmt.Println("Starting Account Stats API")
	HandleRequests()
}
