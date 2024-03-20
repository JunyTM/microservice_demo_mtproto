package main

import (
	"log"
	"ms_gmail/router"
	"net/http"
	"time"
)

func main() {
	log.Printf("Starting micro_gmail: port - %s\n", "8080")
	s := http.Server{
		Addr:    ":8080",
		Handler: router.Router(),

		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 600 * time.Second,
		ReadTimeout:  600 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
