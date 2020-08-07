package main

import (
	"log"
	"net/http"
	"os"
	"star-wars/api"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := ":" + os.Getenv("PORT")
	if port == ":" {
		log.Fatal("PORT must be set")
	}

	s := &http.Server{
		Addr:           port,
		Handler:        api.Config(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Panic("error at listen and serve", s.ListenAndServe())
}
