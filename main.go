package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	LoadConfigFromEnv()
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	server := &http.Server{
		Addr:         "127.0.0.1:" + ServerPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 180 * time.Second,
	}
	log.Println("Go server on http://127.0.0.1:" + ServerPort)
	log.Fatal(server.ListenAndServe())
}
