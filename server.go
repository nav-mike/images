package main

import (
	"fmt"
	"net/http"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	fmt.Println(r.Body)
	fmt.Fprint(w, "[POST] UploadHandler")
}

func main() {
	http.HandleFunc("/upload", UploadHandler)
	http.ListenAndServe(":8080", nil)
}
