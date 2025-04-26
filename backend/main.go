package main

import (
	"log"
	"net/http"

	"github.com/mcrors/ytd/api"
)

func main() {
	server := api.NewServer()

	log.Println("server running on port 8080 ...")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatal(err)
	}
}
