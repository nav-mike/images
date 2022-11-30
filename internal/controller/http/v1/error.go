package v1

import (
	"log"
	"net/http"
)

// errorResponse renders error message to response and logs it
func errorResponse(w http.ResponseWriter, message string, status int, err error) {
	if err != nil {
		log.Printf("Error %s\n", err)
	} else {
		log.Printf("Error %s\n", message)
	}
	http.Error(w, message, status)
}
