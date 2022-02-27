package main

import (
	"log"
	"net/http"

	"github.com/os4ua/browser-crypto-mining/server/internal/server"
)

func main() {
	s, err := server.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting server")
	err = http.ListenAndServe(":8080", s)
	log.Fatal(err)
}
