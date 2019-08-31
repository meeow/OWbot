package accountstatsapi

import (
	"fmt"
	"log"
	"net/http"

	"../accountstats"
)

var (
	testbtag = []string{"brokenglass-4115,orisa"}
)

func homePage(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Fprintf(w, accountstats.ConcurrentGetRawAccountStats(testbtag)[0])
	fmt.Println("Endpoint Hit: homePage")
}

func HandleRequests() {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
