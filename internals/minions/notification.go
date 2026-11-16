package minions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

// NotificationHandler processes outgoing SMS and email tasks.
type NotificationHandler struct {
	log *logrus.Entry
}

// NewNotificationHandler creates a handler for notification tasks.
func NewNotificationHandler(log *logrus.Entry) *NotificationHandler {
	if log == nil {
		log = logrus.NewEntry(logrus.StandardLogger())
	}
	return &NotificationHandler{log: log}
}

// HandleSendSMS processes an SMS delivery task.
// TODO: integrate with BriqSMS (or your provider of choice).
func (h *NotificationHandler) HandleSendSMS(ctx context.Context, t *asynq.Task) error {
	var p SMSPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("unmarshal SMS payload: %w", err)
	}

	h.log.WithField("to", p.To).Info("Sending SMS")

	// ── BriqSMS integration placeholder ──
	// err := briqsms.Send(p.To, p.Message)
	// if err != nil {
	//     return fmt.Errorf("briqsms send: %w", err)
	// }

	h.log.WithField("to", p.To).Info("SMS delivered")
	return nil
}

// HandleSendEmail processes an email delivery task.
// TODO: integrate with Postmark (or your provider of choice).
func (h *NotificationHandler) HandleSendEmail(ctx context.Context, t *asynq.Task) error {
	var p EmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("unmarshal email payload: %w", err)
	}

	h.log.WithFields(logrus.Fields{"to": p.To, "subject": p.Subject}).Info("Sending email")

	// ── Postmark integration placeholder ──
	// err := postmark.Send(p.To, p.Subject, p.Body)
	// if err != nil {
	//     return fmt.Errorf("postmark send: %w", err)
	// }

	h.log.WithField("to", p.To).Info("Email delivered")
	return nil
}
