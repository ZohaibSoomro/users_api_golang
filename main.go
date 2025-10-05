package main

import (
	"log"

	"github.com/zohaibsoomro/users_api_golang/api"
)

func main() {

	api := api.NewApiWithAddress("localhost:8080")
	server := api.RegisterHandlers()
	log.Println("Started server at :8080")
	log.Fatal(server.Run(api.Address))

}
