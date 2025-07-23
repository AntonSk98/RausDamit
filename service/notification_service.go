package service

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"raus-damit/config"
	"strings"
)

type NotificationService struct {
	Config *config.Config
}

func NewNotificationService(config *config.Config) *NotificationService {
	return &NotificationService{config}
}

func (notifier *NotificationService) NotifyCalendarReplacement() error {
	template := notifier.Config.FindTemplateBy(string(REPLACE_CALENDAR_REMINDER))
	emailBody := notifier.buildEmail(template, nil)

	return notifier.sendEmail(emailBody)
}

func (notifier *NotificationService) Notify(notification RubbishCollectionNotification) error {
	if len(notification.RubbishEvents) == 0 {
		return nil
	}

	template := notifier.Config.FindTemplateBy(string(notification.TemplateType))
	templateContent := map[string]string{
		"events": notification.PrettyPrint(),
	}

	emailBody := notifier.buildEmail(template, templateContent)

	return notifier.sendEmail(emailBody)

}

func (notifier *NotificationService) sendEmail(emailBody []byte) error {
	err := smtp.SendMail(
		notifier.Config.Email.SMTPHost+":"+notifier.Config.Email.SMTPPort,
		smtp.PlainAuth("", notifier.Config.Email.From, notifier.Config.Email.Password, notifier.Config.Email.SMTPHost),
		notifier.Config.Email.From,
		notifier.Config.Email.Recepients,
		emailBody,
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Println("✅ Notification has been successfully sent.")
	return nil
}

func (notifier *NotificationService) buildEmail(template config.Template, templateVars map[string]string) []byte {
	subject := template.Subject

	content := template.ResolveContent(templateVars)

	var msg bytes.Buffer

	fromAddress := notifier.Config.Email.From
	toAddresses := strings.Join(notifier.Config.Email.Recepients, ",")

	msg.WriteString("From: Raus Damit <" + fromAddress + ">\r\n")
	msg.WriteString("To: " + toAddresses + "\r\n")
	msg.WriteString("Subject: " + subject + "\r\n")
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	msg.WriteString("Content-Transfer-Encoding: 7bit\r\n")

	// Header/body separator
	msg.WriteString("\r\n")
	msg.WriteString(content)

	return msg.Bytes()
}
