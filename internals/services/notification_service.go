package services

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/tacheraSasi/go-api-starter/internals/tasks"
)

// NotificationService enqueues async notification tasks via asynq.
type NotificationService struct {
	client *asynq.Client
}

// NewNotificationService creates a NotificationService backed by the given asynq client.
func NewNotificationService(client *asynq.Client) *NotificationService {
	return &NotificationService{client: client}
}

// SendSMSAsync enqueues an SMS delivery task on the "critical" queue (max 3 retries).
func (s *NotificationService) SendSMSAsync(to, message string) error {
	payload, err := json.Marshal(tasks.SMSPayload{To: to, Message: message})
	if err != nil {
		return fmt.Errorf("marshal SMS payload: %w", err)
	}

	task := asynq.NewTask(tasks.TypeNotificationSMS, payload)
	_, err = s.client.Enqueue(task,
		asynq.Queue("critical"),
		asynq.MaxRetry(3),
	)
	if err != nil {
		return fmt.Errorf("enqueue SMS task: %w", err)
	}
	return nil
}

// SendEmailAsync enqueues an email delivery task on the "critical" queue (max 3 retries).
func (s *NotificationService) SendEmailAsync(to, subject, body string) error {
	payload, err := json.Marshal(tasks.EmailPayload{To: to, Subject: subject, Body: body})
	if err != nil {
		return fmt.Errorf("marshal email payload: %w", err)
	}

	task := asynq.NewTask(tasks.TypeNotificationEmail, payload)
	_, err = s.client.Enqueue(task,
		asynq.Queue("critical"),
		asynq.MaxRetry(3),
	)
	if err != nil {
		return fmt.Errorf("enqueue email task: %w", err)
	}
	return nil
}
