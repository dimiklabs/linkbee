package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shafikshaon/linkbee/logger"
	"github.com/shafikshaon/linkbee/model"
	"github.com/shafikshaon/linkbee/repository"
	clickSvc "github.com/shafikshaon/linkbee/service/click"
	"github.com/valkey-io/valkey-go/valkeycompat"
)

const bioClickQueueKey = "queue:bio_clicks"

// BioClickWorker consumes bio link click events from the Valkey queue and
// batch-inserts them into PostgreSQL, then updates bio_links.click_count.
type BioClickWorker struct {
	cache                  valkeycompat.Cmdable
	bioLinkClickEventRepo  repository.BioLinkClickEventRepositoryI
	bioRepo                repository.BioRepositoryI
	batchSize              int
	flushInterval          time.Duration
}

func NewBioClickWorker(
	cache valkeycompat.Cmdable,
	bioLinkClickEventRepo repository.BioLinkClickEventRepositoryI,
	bioRepo repository.BioRepositoryI,
	batchSize, flushIntervalSeconds int,
) *BioClickWorker {
	return &BioClickWorker{
		cache:                 cache,
		bioLinkClickEventRepo: bioLinkClickEventRepo,
		bioRepo:               bioRepo,
		batchSize:             batchSize,
		flushInterval:         time.Duration(flushIntervalSeconds) * time.Second,
	}
}

// Start runs the bio click worker in the background. Call this in a goroutine.
func (w *BioClickWorker) Start(ctx context.Context) {
	logger.Info("Bio click worker started",
		zap.Int("batch_size", w.batchSize),
		zap.Duration("flush_interval", w.flushInterval))

	ticker := time.NewTicker(w.flushInterval)
	defer ticker.Stop()

	var batch []*model.BioLinkClickEvent

	for {
		select {
		case <-ctx.Done():
			if len(batch) > 0 {
				w.flush(ctx, batch)
			}
			logger.Info("Bio click worker stopped")
			return

		case <-ticker.C:
			batch = w.drain(ctx, batch)
			if len(batch) > 0 {
				w.flush(ctx, batch)
				batch = batch[:0]
			}
		}
	}
}

func (w *BioClickWorker) drain(ctx context.Context, batch []*model.BioLinkClickEvent) []*model.BioLinkClickEvent {
	for len(batch) < w.batchSize {
		result, err := w.cache.LPop(ctx, bioClickQueueKey).Result()
		if err != nil {
			break
		}
		var payload clickSvc.BioClickPayload
		if err := json.Unmarshal([]byte(result), &payload); err != nil {
			logger.Error("Failed to unmarshal bio click payload", zap.Error(err))
			continue
		}
		batch = append(batch, clickSvc.ToBioLinkClickEvent(payload))
	}
	return batch
}

func (w *BioClickWorker) flush(ctx context.Context, batch []*model.BioLinkClickEvent) {
	if err := w.bioLinkClickEventRepo.BulkCreate(ctx, batch); err != nil {
		logger.Error("Failed to flush bio click events to database",
			zap.Error(err),
			zap.Int("batch_size", len(batch)))
		return
	}
	logger.Debug("Flushed bio click events to database", zap.Int("count", len(batch)))

	// Tally clicks per bio link and update click_count atomically.
	counts := make(map[uuid.UUID]int64, len(batch))
	for _, e := range batch {
		counts[e.BioLinkID]++
	}
	for bioLinkID, delta := range counts {
		if err := w.bioRepo.AddBioLinkClickCount(ctx, bioLinkID, delta); err != nil {
			logger.Error("Failed to update bio link click_count",
				zap.String("bio_link_id", bioLinkID.String()),
				zap.Int64("delta", delta),
				zap.Error(err))
		}
	}
}
