package main

import (
	"log"
	"net/http"
)

func main() {
	RegisterRoutes()

	log.Println("Dashboard running on: http://10.8.0.1:8080")
	log.Fatal(http.ListenAndServe("10.8.0.1:8080", nil))
}