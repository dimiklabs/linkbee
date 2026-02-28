package db

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/shafikshaon/linkbee/logger"
)

// SetupTimescale enables the TimescaleDB extension and configures the
// click_events hypertable, compression, and a continuous aggregate.
// It is idempotent and safe to call on every server start.
func SetupTimescale(db *gorm.DB) error {
	logger.Info("Setting up TimescaleDB")

	// 1. Enable extension (must precede all TimescaleDB API calls).
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE`).Error; err != nil {
		logger.Error("Failed to enable timescaledb extension", zap.Error(err))
		return err
	}

	// 2. Convert click_events to a hypertable partitioned by clicked_at.
	//    chunk_time_interval = 1 day keeps individual chunks small so
	//    time-range queries touch only the relevant day partitions.
	//    migrate_data = TRUE is a no-op on an empty table and handles
	//    existing rows when upgrading a running deployment.
	if err := db.Exec(`
		SELECT create_hypertable(
			'click_events', 'clicked_at',
			chunk_time_interval => INTERVAL '1 day',
			if_not_exists       => TRUE,
			migrate_data        => TRUE
		)
	`).Error; err != nil {
		logger.Error("Failed to create hypertable for click_events", zap.Error(err))
		return err
	}

	// 3. Enable column-store compression on the hypertable.
	//    Segmenting by link_id groups each link's events together within a
	//    compressed chunk, so per-link analytic scans decompress the minimum
	//    amount of data.
	if err := db.Exec(`
		ALTER TABLE click_events SET (
			timescaledb.compress,
			timescaledb.compress_segmentby = 'link_id',
			timescaledb.compress_orderby   = 'clicked_at DESC'
		)
	`).Error; err != nil {
		logger.Warn("Failed to configure compression for click_events (may already be set)", zap.Error(err))
	}

	// 4. Automatically compress chunks that are older than 7 days.
	if err := db.Exec(`
		SELECT add_compression_policy(
			'click_events', INTERVAL '7 days', if_not_exists => TRUE
		)
	`).Error; err != nil {
		logger.Warn("Failed to add compression policy for click_events", zap.Error(err))
	}

	// 5. Hourly continuous aggregate — pre-computes total and unique click
	//    counts per link per hour so the common "clicks over time" dashboard
	//    query reads from the materialized view rather than raw events.
	if err := db.Exec(`
		CREATE MATERIALIZED VIEW IF NOT EXISTS click_stats_hourly
		WITH (timescaledb.continuous) AS
		SELECT
			link_id,
			time_bucket('1 hour', clicked_at) AS bucket,
			COUNT(*)                           AS total_clicks,
			COUNT(DISTINCT ip_hash)            AS unique_clicks
		FROM click_events
		GROUP BY link_id, bucket
		WITH NO DATA
	`).Error; err != nil {
		logger.Warn("Failed to create click_stats_hourly continuous aggregate", zap.Error(err))
	}

	// 6. Keep the continuous aggregate up to date: refresh every hour,
	//    covering events from 3 hours ago up to 1 hour ago (1-hour lag
	//    ensures the most recent open chunk is not yet being refreshed).
	if err := db.Exec(`
		SELECT add_continuous_aggregate_policy(
			'click_stats_hourly',
			start_offset      => INTERVAL '3 hours',
			end_offset        => INTERVAL '1 hour',
			schedule_interval => INTERVAL '1 hour',
			if_not_exists     => TRUE
		)
	`).Error; err != nil {
		logger.Warn("Failed to add refresh policy for click_stats_hourly", zap.Error(err))
	}

	logger.Info("TimescaleDB setup completed")
	return nil
}
