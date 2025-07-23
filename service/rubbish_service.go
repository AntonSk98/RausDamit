package service

import (
	"log"
	"raus-damit/config"
	"raus-damit/repository"
)

type RubbishEventService struct {
	rubbishEventRepository *repository.RubbishEventRepository
	config                 *config.Config
	notificationService    *NotificationService
}

func NewRubbishEventService(notificationService *NotificationService,
	rubbishEventRepository *repository.RubbishEventRepository,
	config *config.Config) *RubbishEventService {

	return &RubbishEventService{rubbishEventRepository, config, notificationService}
}

func (service *RubbishEventService) NotifyDailyRubbishCollection() error {
	now := DateNow(service.config.Timezone)
	tomorrow := PlusDays(now, 1)
	events := service.rubbishEventRepository.Find(now, tomorrow)

	if len(events) == 0 {
		log.Println("📭 No rubbish collection scheduled for tomorrow.")
		return nil
	}

	notification := NewRubbishCollectionNotification(DAILY_RUBBISH_REMINDER, events)
	return service.notificationService.Notify(notification)
}

func (service *RubbishEventService) NotifyWeeklyRubbishCollection() error {
	now := DateNow(service.config.Timezone)
	nextWeek := PlusDays(now, 7)
	events := service.rubbishEventRepository.Find(now, nextWeek)

	if len(events) == 0 {
		log.Println("📭 No rubbish collection scheduled in the next 7 days.")
		return nil
	}

	notification := NewRubbishCollectionNotification(WEEKLY_RUBBISH_REMINDER, events)
	return service.notificationService.Notify(notification)
}

func (service *RubbishEventService) NotifyReplaceCalendar() error {
	return service.notificationService.NotifyCalendarReplacement()
}
