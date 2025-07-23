package main

import (
	"log"
	"os"
	"raus-damit/config"
	"raus-damit/repository"
	"raus-damit/service"
	"time"

	"github.com/avast/retry-go"
	"github.com/robfig/cron/v3"
)

func main() {
	configuration := config.LoadConfig(os.Getenv("CONFIG_PATH"))
	rubbishEventRepository := repository.NewRubbishEventRepository(configuration)
	notificationService := service.NewNotificationService(configuration)
	rubbishEventService := service.NewRubbishEventService(notificationService, rubbishEventRepository, configuration)

	cronScheduler := cron.New(cron.WithSeconds())

	// Every minute
	cronScheduler.AddFunc("* * * * * *", func() {
		retriable("DailyMorningNotifier", func() error { return rubbishEventService.NotifyDailyRubbishCollection() })
	})

	// Every day at 1:00 PM
	cronScheduler.AddFunc("0 0 13 * * *", func() {
		retriable("DailyMorningNotifier", func() error { return rubbishEventService.NotifyDailyRubbishCollection() })
	})

	// Every day at 6:00 PM
	cronScheduler.AddFunc("0 0 18 * * *", func() {
		retriable("EveningMorningNotifier", func() error { return rubbishEventService.NotifyDailyRubbishCollection() })
	})

	// Every Sunday at 8:00 AM
	cronScheduler.AddFunc("0 0 8 * * 0", func() {
		retriable("WeeklyNotifier", func() error { return rubbishEventService.NotifyWeeklyRubbishCollection() })
	})

	// Every year on December 30th at 10:00 AM
	cronScheduler.AddFunc("0 0 10 30 12 *", func() {
		retriable("YearlyCalendarReplacementNotifier", func() error {
			return rubbishEventService.NotifyReplaceCalendar()
		})
	})

	cronScheduler.Start()

	// Keep the program alive
	select {}
}

func retriable(taskId string, task func() error) error {
	return retry.Do(
		task,
		retry.Attempts(20),
		retry.Delay(5*time.Minute),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("Retry #%d of '%s()' due to error: %v", n+1, taskId, err)
		}),
	)
}
