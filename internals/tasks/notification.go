package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

// HandleSMSTask processes an SMS delivery task.
// TODO: integrate with BriqSMS (or your provider of choice).
func HandleSMSTask(ctx context.Context, t *asynq.Task) error {
	var p SMSPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("unmarshal SMS payload: %w", err)
	}

	log.Printf("[SMS] sending to=%s msg=%s", p.To, p.Message)

	// ── BriqSMS integration placeholder ──
	// err := briqsms.Send(p.To, p.Message)
	// if err != nil {
	//     return fmt.Errorf("briqsms send: %w", err)
	// }

	log.Printf("[SMS] delivered to=%s", p.To)
	return nil
}

// HandleEmailTask processes an email delivery task.
// TODO: integrate with Postmark (or your provider of choice).
func HandleEmailTask(ctx context.Context, t *asynq.Task) error {
	var p EmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("unmarshal email payload: %w", err)
	}

	log.Printf("[Email] sending to=%s subject=%q", p.To, p.Subject)

	// ── Postmark integration placeholder ──
	// err := postmark.Send(p.To, p.Subject, p.Body)
	// if err != nil {
	//     return fmt.Errorf("postmark send: %w", err)
	// }

	log.Printf("[Email] delivered to=%s", p.To)
	return nil
}
