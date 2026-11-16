package minions

// Task type constants for the asynq queue.
// Naming convention: "domain:action".
//
// Add new task types here and register their handlers in cmd/minion. Schedule
// recurring tasks in cmd/overseer. Keep business logic in internals/services or
// internals/minions handlers so the binaries stay thin.
const (
	// Notifications
	TaskSendEmail = "notification:email"
	TaskSendSMS   = "notification:sms"

	// Maintenance (placeholder used to demonstrate the overseer -> minion flow)
	TaskHeartbeat = "maintenance:heartbeat"
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
