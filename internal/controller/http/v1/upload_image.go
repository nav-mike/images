package v1

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nav-mike/images/config"
	"github.com/nav-mike/images/internal/entity"
)

func PostUploadImageHandler(config *config.Config, writer ImageWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Upload image")
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var file entity.UploadImageDTO

		err := json.NewDecoder(r.Body).Decode(&file)
		if err != nil {
			errorResponse(w, "Bad request", http.StatusBadRequest, err)
			return
		}

		result, err := writer.SaveImage(file)
		if err != nil {
			errorResponse(w, "Internal server error", http.StatusInternalServerError, err)
			return
		}

		for key, value := range result {
			result[key] = imageUrl(config, value)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func imageUrl(config *config.Config, filename string) string {
	return config.ServerHostUrl + "/images/" + filename
}
