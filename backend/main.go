package main

import (
	"log"
	"net/http"

	"petwell-merchant-backend/internal"
)

func main() {
	store := internal.NewStore()
	store.SeedDemoData()
	server := internal.NewServer(store)

	addr := ":8090"
	log.Printf("PetWell merchant backend running at http://localhost%s", addr)
	if err := http.ListenAndServe(addr, server.Routes()); err != nil {
		log.Fatal(err)
	}
}
