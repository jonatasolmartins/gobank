package main

import (
	"flag"
	"log"
)

func main() {

	seed := flag.Bool("seed", false, "seed the database")
	flag.Parse()

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		if err := store.Seed(); err != nil {
			log.Fatal(err)
		}
	}
	server := NewAPIServer(":3001", store)
	server.Run()
}
