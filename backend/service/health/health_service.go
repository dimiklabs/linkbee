package health

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/valkey-io/valkey-go/valkeycompat"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/response"
)

const (
	StatusUp   = "up"
	StatusDown = "down"
	Version    = "1.0.0"
)

type HealthServiceI interface {
	Check(ctx context.Context) *response.HealthResponse
}

type HealthService struct {
	masterDB  *gorm.DB
	replicaDB *gorm.DB
	cache     valkeycompat.Cmdable
	env       string
	startTime time.Time
}

func NewHealthService(masterDB *gorm.DB, replicaDB *gorm.DB, cache valkeycompat.Cmdable, env string) HealthServiceI {
	return &HealthService{
		masterDB:  masterDB,
		replicaDB: replicaDB,
		cache:     cache,
		env:       env,
		startTime: time.Now(),
	}
}

func (s *HealthService) Check(ctx context.Context) *response.HealthResponse {
	checks := make(map[string]response.ComponentCheck)

	checks["database_master"] = s.checkDatabase(ctx, s.masterDB, "master")
	checks["database_replica"] = s.checkDatabase(ctx, s.replicaDB, "replica")
	checks["valkey"] = s.checkValkey(ctx)

	overallStatus := StatusUp
	for name, check := range checks {
		if check.Status == StatusDown {
			overallStatus = StatusDown
			logger.WarnCtx(ctx, "Health check component is down",
				zap.String("component", name),
				zap.String("error", check.Error))
		}
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	resp := &response.HealthResponse{
		Status: overallStatus,
		Uptime: time.Since(s.startTime).Round(time.Second).String(),
		Checks: checks,
		System: response.SystemInfo{
			Version:     Version,
			Env:         s.env,
			GoRoutines:  runtime.NumGoroutine(),
			MemoryAlloc: formatBytes(memStats.Alloc),
			MemorySys:   formatBytes(memStats.Sys),
		},
	}

	logger.DebugCtx(ctx, "Health check completed",
		zap.String("status", overallStatus),
		zap.Int("go_routines", runtime.NumGoroutine()))

	return resp
}

func (s *HealthService) checkDatabase(ctx context.Context, db *gorm.DB, dbType string) response.ComponentCheck {
	if db == nil {
		return response.ComponentCheck{
			Status: StatusDown,
			Error:  "database connection not initialized",
		}
	}

	start := time.Now()
	sqlDB, err := db.DB()
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get database instance for health check",
			zap.String("type", dbType),
			zap.Error(err))
		return response.ComponentCheck{
			Status: StatusDown,
			Error:  err.Error(),
		}
	}

	err = sqlDB.PingContext(ctx)
	latency := time.Since(start)

	if err != nil {
		logger.ErrorCtx(ctx, "Database ping failed",
			zap.String("type", dbType),
			zap.Duration("latency", latency),
			zap.Error(err))
		return response.ComponentCheck{
			Status:  StatusDown,
			Latency: latency.String(),
			Error:   err.Error(),
		}
	}

	return response.ComponentCheck{
		Status:  StatusUp,
		Latency: latency.String(),
	}
}

func (s *HealthService) checkValkey(ctx context.Context) response.ComponentCheck {
	if s.cache == nil {
		return response.ComponentCheck{
			Status: StatusDown,
			Error:  "valkey connection not initialized",
		}
	}

	start := time.Now()
	err := s.cache.Ping(ctx).Err()
	latency := time.Since(start)

	if err != nil {
		logger.ErrorCtx(ctx, "Valkey ping failed",
			zap.Duration("latency", latency),
			zap.Error(err))
		return response.ComponentCheck{
			Status:  StatusDown,
			Latency: latency.String(),
			Error:   err.Error(),
		}
	}

	return response.ComponentCheck{
		Status:  StatusUp,
		Latency: latency.String(),
	}
}

func formatBytes(bytes uint64) string {
	const (
		mb = 1024 * 1024
	)
	return fmt.Sprintf("%.2f MB", float64(bytes)/float64(mb))
}
