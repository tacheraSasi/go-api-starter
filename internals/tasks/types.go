package tasks

// Task type constants used by both the API (enqueue) and worker (dequeue).
const (
	TypeNotificationSMS   = "notification:sms"
	TypeNotificationEmail = "notification:email"
)

// SMSPayload is the JSON payload for an SMS task.
type SMSPayload struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

// EmailPayload is the JSON payload for an email task.
type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	HTML    string `json:"html,omitempty"`
}
