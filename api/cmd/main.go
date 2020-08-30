package main

import (
	"log"
	"net/http"
	"star-wars/api"
	"star-wars/env"
	"time"
)

func main() {
	port := ":" + env.Vars.Api.Port

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
