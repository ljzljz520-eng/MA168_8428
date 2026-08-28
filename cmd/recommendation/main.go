package main

import (
	"bookstore/recommendation/internal/api"
	"bookstore/recommendation/internal/store"
	"bookstore/recommendation/internal/workflow"
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	dbPath := flag.String("db", "recommendations.db", "path to bbolt database")
	addr := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()
	database, err := store.Open(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	service := workflow.New(database)
	server := api.New(service, log.New(os.Stdout, "recommendation ", log.LstdFlags))
	log.Printf("bookstore recommendation service listening on %s", *addr)
	if err := http.ListenAndServe(*addr, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
