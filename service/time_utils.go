package service

import (
	"log"
	"time"
)

func Location(timezone string) *time.Location {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}

	return loc
}

func DateNow(timezone string) time.Time {
	return time.Now().In(Location(timezone))
}

func PlusDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}
