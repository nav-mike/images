package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UploadImageDTO struct {
	File string
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	var file UploadImageDTO

	err := json.NewDecoder(r.Body).Decode(&file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	result := map[string]string{"result": "success"}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/upload", UploadHandler)
	http.ListenAndServe(":8080", nil)
}
