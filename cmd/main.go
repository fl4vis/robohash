package main

import (
	"log"

	"github.com/fl4vis/robohash/internal/app/robohash"
)

func main() {
	log.SetFlags(0)
	if err := robohash.Run(); err != nil {
		log.Fatal(err)
	}
}
