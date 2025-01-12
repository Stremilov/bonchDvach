package main

import (
	"bonchDvach"
	"bonchDvach/pkg/handlers"
	_ "bonchDvach/pkg/handlers"
	"log"
)

func main() {
	server := new(bonchDvach.Server)

	if err := server.Run("8000", handlers.InitRoutesAndDB()); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}
