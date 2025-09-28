package main

import (
	"log"
	"net/http"

	"github.com/zohaibsoomro/golangpractice/api"
)

func main() {

	api := api.NewApi("localhost:8080")
	api.RegisterHandlers()
	log.Println("Started server at :8080")
	log.Fatal(http.ListenAndServe(api.Address, nil))
	log.Println("Starting stopped.")

}
