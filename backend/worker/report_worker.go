package worker

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/shafikshaon/shortlink/logger"
	reportSvc "github.com/shafikshaon/shortlink/service/reporting"
)

// ReportWorker triggers scheduled analytics report delivery every minute.
type ReportWorker struct {
	svc      reportSvc.ReportingServiceI
	interval time.Duration
}

func NewReportWorker(svc reportSvc.ReportingServiceI) *ReportWorker {
	return &ReportWorker{
		svc:      svc,
		interval: time.Minute,
	}
}

// Start runs the report worker in the background. Call this in a goroutine.
func (w *ReportWorker) Start(ctx context.Context) {
	logger.Info("Report worker started", zap.Duration("interval", w.interval))

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Report worker stopped")
			return
		case <-ticker.C:
			w.svc.ProcessDueReports(ctx)
		}
	}
}
