package initapi

import (
	"fmt"

	"../accountstatsapi"
)

// StartAllEndpoints starts the REST API service
func StartAllEndpoints() {
	fmt.Println("Starting Account Stats API")
	go accountstatsapi.HandleRequests()
}
