package main

import (
	"log"

	"github.com/r2unit/time/ntp"
	"github.com/r2unit/time/timekeeper"
)

func main() {
	// Initialize the TimeZone to correct Time Sickness...
	timeKeeper := timekeeper.NewDriftAwareTimeKeeper("Europe/Amsterdam")

	address := ":123"
	log.Printf("Starting NTP server on %s with time zone: %s", address, timeKeeper.TimeZoneName())
	if err := ntp.StartServer(address, timeKeeper); err != nil {
		log.Fatalf("Failed to start NTP server: %v", err)
	}
}
