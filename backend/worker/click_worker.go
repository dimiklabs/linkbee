package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/repository"
	clickSvc "github.com/shafikshaon/shortlink/service/click"
	"github.com/valkey-io/valkey-go/valkeycompat"
)

const (
	clickQueueKey = "queue:clicks"
)

// ClickWorker consumes click events from the Valkey queue and batch-inserts them into PostgreSQL.
type ClickWorker struct {
	cache          valkeycompat.Cmdable
	clickEventRepo repository.ClickEventRepositoryI
	linkRepo       repository.LinkRepositoryI
	batchSize      int
	flushInterval  time.Duration
}

func NewClickWorker(cache valkeycompat.Cmdable, clickEventRepo repository.ClickEventRepositoryI, linkRepo repository.LinkRepositoryI, batchSize, flushIntervalSeconds int) *ClickWorker {
	return &ClickWorker{
		cache:          cache,
		clickEventRepo: clickEventRepo,
		linkRepo:       linkRepo,
		batchSize:      batchSize,
		flushInterval:  time.Duration(flushIntervalSeconds) * time.Second,
	}
}

// Start runs the click worker in the background. Call this in a goroutine.
func (w *ClickWorker) Start(ctx context.Context) {
	logger.Info("Click worker started",
		zap.Int("batch_size", w.batchSize),
		zap.Duration("flush_interval", w.flushInterval))

	ticker := time.NewTicker(w.flushInterval)
	defer ticker.Stop()

	var batch []*model.ClickEvent

	for {
		select {
		case <-ctx.Done():
			// Flush remaining events on shutdown
			if len(batch) > 0 {
				w.flush(ctx, batch)
			}
			logger.Info("Click worker stopped")
			return

		case <-ticker.C:
			// Drain the queue up to batchSize
			batch = w.drain(ctx, batch)
			if len(batch) > 0 {
				w.flush(ctx, batch)
				batch = batch[:0]
			}
		}
	}
}

func (w *ClickWorker) drain(ctx context.Context, batch []*model.ClickEvent) []*model.ClickEvent {
	for len(batch) < w.batchSize {
		result, err := w.cache.LPop(ctx, clickQueueKey).Result()
		if err != nil {
			// Queue is empty or error
			break
		}

		var payload clickSvc.ClickPayload
		if err := json.Unmarshal([]byte(result), &payload); err != nil {
			logger.Error("Failed to unmarshal click payload", zap.Error(err))
			continue
		}

		batch = append(batch, clickSvc.ToClickEvent(payload))
	}
	return batch
}

func (w *ClickWorker) flush(ctx context.Context, batch []*model.ClickEvent) {
	if err := w.clickEventRepo.BulkCreate(ctx, batch); err != nil {
		logger.Error("Failed to flush click events to database",
			zap.Error(err),
			zap.Int("batch_size", len(batch)))
		return
	}
	logger.Debug("Flushed click events to database", zap.Int("count", len(batch)))

	// Tally clicks per link and update click_count on the links table.
	counts := make(map[uuid.UUID]int64, len(batch))
	for _, e := range batch {
		counts[e.LinkID]++
	}
	for linkID, delta := range counts {
		if err := w.linkRepo.AddClickCount(ctx, linkID, delta); err != nil {
			logger.Error("Failed to update link click_count",
				zap.String("link_id", linkID.String()),
				zap.Int64("delta", delta),
				zap.Error(err))
		}
	}
}
