package main

import (
	"bonchDvach"
	_ "bonchDvach/docs"
	"bonchDvach/pkg/router"
	"log"
)

// @title           BonchDvach API
// @version         1.0
// @description     API сервиса BonchDvach

func main() {
	server := new(bonchDvach.Server)

	if err := server.Run("8000", router.InitRoutesAndDB()); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}
