package service

import (
	"fmt"
	"raus-damit/repository"
	"strings"
	"time"
)

// TemplateType defines the type for template keys
type TemplateType string

const (
	// DAILY_RUBBISH_REMINDER reminds one day before the rubbish collection
	DAILY_RUBBISH_REMINDER TemplateType = "daily_rubbish_reminder"

	// WEEKLY_RUBBISH_REMINDER sends a calendar with all collection dates for the week
	WEEKLY_RUBBISH_REMINDER TemplateType = "weekly_rubbish_reminder"

	// REPLACE_CALENDAR_REMINDER notifies to replace calendar for the new year
	REPLACE_CALENDAR_REMINDER TemplateType = "replace_calendar_reminder"
)

type RubbishCollectionNotification struct {
	TemplateType  TemplateType
	RubbishEvents []RubbishEventInfo
}

type RubbishEventInfo struct {
	Date        time.Time
	RubbishType string
}

func NewRubbishCollectionNotification(templateType TemplateType, rubishEvents []repository.RubbishEvent) RubbishCollectionNotification {
	var rubbishEventInfos []RubbishEventInfo

	for _, event := range rubishEvents {
		rubbishEventInfos = append(rubbishEventInfos, RubbishEventInfo{event.Date, event.Fraction})
	}

	return RubbishCollectionNotification{
		TemplateType:  templateType,
		RubbishEvents: rubbishEventInfos,
	}
}

func (notification *RubbishCollectionNotification) PrettyPrint() string {
	lines := make([]string, 0, len(notification.RubbishEvents))
	for _, event := range notification.RubbishEvents {
		line := fmt.Sprintf("🗑️ %s | %s", event.Date.Format(time.DateOnly), event.RubbishType)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
