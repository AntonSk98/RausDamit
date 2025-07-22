package service

import (
	"log"
	"time"
)

func DateNow(timezone string) time.Time {
	location, err := time.LoadLocation(timezone)

	if err != nil {
		log.Fatal(err)
	}

	now := time.Now().In(location)

	return now
}

func PlusDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}
