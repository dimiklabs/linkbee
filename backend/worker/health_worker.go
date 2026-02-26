package worker

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/repository"
	"github.com/shafikshaon/shortlink/response"
)

const (
	healthCheckInterval = 30 * time.Minute
	healthCheckBatch    = 50
	healthStaleDuration = 24 * time.Hour
	healthProbeTimeout  = 10 * time.Second
)

// HealthWorker periodically checks the reachability of active destination URLs.
type HealthWorker struct {
	linkRepo   repository.LinkRepositoryI
	httpClient *http.Client
}

func NewHealthWorker(linkRepo repository.LinkRepositoryI) *HealthWorker {
	return &HealthWorker{
		linkRepo: linkRepo,
		httpClient: &http.Client{
			Timeout: healthProbeTimeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
	}
}

// Start runs the health worker in the background. Call this in a goroutine.
func (w *HealthWorker) Start(ctx context.Context) {
	logger.Info("Health worker started",
		zap.Duration("interval", healthCheckInterval),
		zap.Int("batch", healthCheckBatch))

	// Run once immediately at startup, then on ticker
	w.runBatch(ctx)

	ticker := time.NewTicker(healthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Health worker stopped")
			return
		case <-ticker.C:
			w.runBatch(ctx)
		}
	}
}

func (w *HealthWorker) runBatch(ctx context.Context) {
	staleBefore := time.Now().Add(-healthStaleDuration)
	links, err := w.linkRepo.GetLinksForHealthCheck(ctx, staleBefore, healthCheckBatch)
	if err != nil {
		logger.Error("Health worker: failed to fetch links", zap.Error(err))
		return
	}

	if len(links) == 0 {
		return
	}

	logger.Info("Health worker: checking links", zap.Int("count", len(links)))

	for _, link := range links {
		select {
		case <-ctx.Done():
			return
		default:
		}

		status, statusCode := w.probe(ctx, link.DestinationURL)
		if updateErr := w.linkRepo.UpdateHealthStatus(ctx, link.ID, status, statusCode); updateErr != nil {
			logger.Error("Health worker: failed to update health status",
				zap.String("link_id", link.ID.String()),
				zap.Error(updateErr))
		} else {
			logger.Debug("Health worker: checked link",
				zap.String("link_id", link.ID.String()),
				zap.String("status", status),
				zap.Int("status_code", statusCode))
		}
	}
}

func (w *HealthWorker) probe(ctx context.Context, rawURL string) (status string, statusCode int) {
	probeCtx, cancel := context.WithTimeout(ctx, healthProbeTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(probeCtx, http.MethodHead, rawURL, nil)
	if err != nil {
		return response.HealthStatusError, 0
	}
	req.Header.Set("User-Agent", "shortlink-health-monitor/1.0")

	resp, err := w.httpClient.Do(req)
	if err != nil {
		if probeCtx.Err() != nil {
			return response.HealthStatusTimeout, 0
		}
		return response.HealthStatusError, 0
	}
	resp.Body.Close()

	// Retry with GET if server rejects HEAD
	if resp.StatusCode == http.StatusMethodNotAllowed || resp.StatusCode == http.StatusNotImplemented {
		getReq, gErr := http.NewRequestWithContext(probeCtx, http.MethodGet, rawURL, nil)
		if gErr == nil {
			getReq.Header.Set("User-Agent", "shortlink-health-monitor/1.0")
			if getResp, gErr2 := w.httpClient.Do(getReq); gErr2 == nil {
				resp.StatusCode = getResp.StatusCode
				getResp.Body.Close()
			}
		}
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return response.HealthStatusHealthy, resp.StatusCode
	}
	return response.HealthStatusUnhealthy, resp.StatusCode
}
