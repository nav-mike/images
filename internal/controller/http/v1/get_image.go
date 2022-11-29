package v1

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nav-mike/images/config"
)

func GetImageHandler(config *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Getting image request %s \n", r.URL.Path)
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userId := getUserId(r)
		if userId == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		fullPath := getFullImagePath(config, userId, r.URL.Path)
		if _, err := os.Stat(fullPath); errors.Is(err, os.ErrNotExist) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		http.ServeFile(w, r, fullPath)
	}
}

func getUserId(r *http.Request) string {
	return r.Header.Get("X-Custom-Auth-Token")
}

func getFullImagePath(config *config.Config, userId, path string) string {
	return config.ImagesDir + "/" + userId + "/" + strings.Replace(path, "/images/", "", 1)
}
