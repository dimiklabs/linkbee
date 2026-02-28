package worker

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/shafikshaon/linkbee/logger"
	expirySvc "github.com/shafikshaon/linkbee/service/expiry"
)

// ExpiryWorker sends email notifications for links that are expiring within 3 days.
type ExpiryWorker struct {
	svc      expirySvc.ExpiryServiceI
	interval time.Duration
}

func NewExpiryWorker(svc expirySvc.ExpiryServiceI) *ExpiryWorker {
	return &ExpiryWorker{
		svc:      svc,
		interval: time.Hour,
	}
}

// Start runs the expiry notification worker. Call this in a goroutine.
func (w *ExpiryWorker) Start(ctx context.Context) {
	logger.Info("Expiry notification worker started", zap.Duration("interval", w.interval))

	// Run once immediately on startup so we don't wait a full hour on boot.
	w.svc.ProcessExpiringLinks(ctx)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Expiry notification worker stopped")
			return
		case <-ticker.C:
			w.svc.ProcessExpiringLinks(ctx)
		}
	}
}
