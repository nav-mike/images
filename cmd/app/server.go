package main

import (
	"log"

	"github.com/nav-mike/images/config"
	v1 "github.com/nav-mike/images/internal/controller/http/v1"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		return
	}

	log.Println("Starting server on port 8080")

	v1.NewRouter(conf)
}
