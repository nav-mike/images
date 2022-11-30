package v1

import (
	"log"
	"net/http"
)

func GetImageHandler(reader ImageReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Getting image request %s \n", r.URL.Path)
		if r.Method != "GET" {
			errorResponse(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
			return
		}

		userId := getUserId(r)
		if userId == "" {
			errorResponse(w, "Unauthorized", http.StatusUnauthorized, nil)
			return
		}

		fullPath, err := reader.GetStaticImagePath(userId, r.URL.Path)
		if err != nil {
			errorResponse(w, "Not found", http.StatusNotFound, err)
			return
		}

		http.ServeFile(w, r, fullPath)
	}
}

func getUserId(r *http.Request) string {
	return r.Header.Get("X-Custom-Auth-Token") // getting fake user from header. In real life it should be from session or JWT
}
