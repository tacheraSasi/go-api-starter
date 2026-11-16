package minions

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

// MaintenanceHandler is the place to put recurring housekeeping jobs such as
// deleting expired rows, pruning logs, or rotating exports.
type MaintenanceHandler struct {
	log *logrus.Entry
}

// NewMaintenanceHandler creates a handler for maintenance tasks.
func NewMaintenanceHandler(log *logrus.Entry) *MaintenanceHandler {
	if log == nil {
		log = logrus.NewEntry(logrus.StandardLogger())
	}
	return &MaintenanceHandler{log: log}
}

// HandleHeartbeat is a placeholder recurring task. It proves the
// overseer -> Redis -> minion path without implying any business rules.
// Replace or extend it with real housekeeping jobs.
func (h *MaintenanceHandler) HandleHeartbeat(ctx context.Context, t *asynq.Task) error {
	h.log.Info("Heartbeat received from overseer")
	return nil
}
