package main

import (
	"log"
	"net/http"
	"os"

	"petwell-merchant-backend/internal"
)

func main() {
	storeType := os.Getenv("PETWELL_STORE")
	if storeType == "" {
		storeType = "sqlite"
	}

	var (
		store *internal.Store
		err   error
	)
	if storeType == "memory" {
		store = internal.NewStore()
		log.Printf("using in-memory store")
	} else {
		dbPath := os.Getenv("PETWELL_DB_PATH")
		if dbPath == "" {
			dbPath = "./db/petwell_merchant.db"
		}
		store, err = internal.NewSQLiteStore(dbPath)
		if err != nil {
			log.Fatalf("init sqlite store: %v", err)
		}
		defer func() {
			if closeErr := store.Close(); closeErr != nil {
				log.Printf("close store: %v", closeErr)
			}
		}()
		log.Printf("using sqlite store at %s", dbPath)
	}

	store.SeedDemoData()
	server := internal.NewServer(store)

	addr := ":8090"
	log.Printf("PetWell merchant backend running at http://localhost%s", addr)
	if err := http.ListenAndServe(addr, server.Routes()); err != nil {
		log.Fatal(err)
	}
}
