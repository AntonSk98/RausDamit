package repository

import (
	"errors"
	"regexp"
	"time"
)

type NewRubbishEventCommand struct {
	City     string
	District string
	Location string
	Fraction string
	Date     string
}

type RubbishEvent struct {
	City     string
	District string
	Location string
	Fraction string
	Date     time.Time
}

func newRubbishEvent(command NewRubbishEventCommand) (*RubbishEvent, error) {
	re := regexp.MustCompile(`\d{2}\.\d{2}\.\d{4}`)
	dateStr := re.FindString(command.Date)

	if dateStr == "" {
		return nil, errors.New("date not found in input")
	}

	formattedDate, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		return nil, err
	}

	rubbishEvent := &RubbishEvent{
		City:     command.City,
		District: command.District,
		Location: command.Location,
		Fraction: command.Fraction,
		Date:     formattedDate,
	}

	return rubbishEvent, nil
}

func (event *RubbishEvent) isBetween(from time.Time, to time.Time) bool {
	return !event.Date.Before(from) && !event.Date.After(to)
}
