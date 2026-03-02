package main

import (
	"log"
	"net/http"

	"petwell-merchant-backend/internal"
)

func main() {
	cfg := internal.LoadConfig()

	var (
		store *internal.Store
		err   error
	)
	if cfg.StoreType == "memory" {
		store = internal.NewStore()
		log.Printf("using in-memory store env=%s addr=%s", cfg.Environment, cfg.Addr)
	} else {
		store, err = internal.NewSQLiteStore(cfg.DBPath)
		if err != nil {
			log.Fatalf("init sqlite store: %v", err)
		}
		defer func() {
			if closeErr := store.Close(); closeErr != nil {
				log.Printf("close store: %v", closeErr)
			}
		}()
		log.Printf("using sqlite store env=%s addr=%s db=%s", cfg.Environment, cfg.Addr, cfg.DBPath)
	}

	store.SeedDemoData()
	server := internal.NewServer(store)

	log.Printf("PetWell merchant backend running at http://localhost%s", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, server.Routes()); err != nil {
		log.Fatal(err)
	}
}
