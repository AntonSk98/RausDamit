package repository

import (
	"encoding/csv"
	"log"
	"net/http"
	"raus-damit/config"
	"time"
)

type RubbishEventRepository struct {
	Events []RubbishEvent
}

func NewRubbishEventRepository(config *config.Config) *RubbishEventRepository {
	records := readRubbishCalendar(config.Calendar)
	rubbishEvents := parseRubbishEvents(records)
	return &RubbishEventRepository{
		Events: rubbishEvents,
	}
}

func (re *RubbishEventRepository) Find(from, to time.Time) []RubbishEvent {
	if len(re.Events) == 0 {
		return []RubbishEvent{}
	}

	var inRange []RubbishEvent
	for _, event := range re.Events {
		if event.isBetween(from, to) {
			inRange = append(inRange, event)
		}
	}

	return inRange
}

func readRubbishCalendar(url string) [][]string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch file: %s", resp.Status)
	}

	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return records
}

func parseRubbishEvents(records [][]string) []RubbishEvent {
	var rubbishEvents []RubbishEvent
	for i, row := range records {
		if i == 0 {
			// skip header row
			continue
		}
		if len(row) < 5 {
			continue
		}

		newRubbishEventCommand := NewRubbishEventCommand{
			City:     row[0],
			District: row[1],
			Location: row[2],
			Fraction: row[3],
			Date:     row[4],
		}

		rc, err := newRubbishEvent(newRubbishEventCommand)

		if err != nil {
			log.Fatal(err)
		}

		rubbishEvents = append(rubbishEvents, *rc)
	}

	return rubbishEvents
}
