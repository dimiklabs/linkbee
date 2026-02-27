package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/model"
)

type AnalyticsReportRepositoryI interface {
	Create(ctx context.Context, report *model.AnalyticsReport) error
	GetByID(ctx context.Context, id, userID uuid.UUID) (*model.AnalyticsReport, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*model.AnalyticsReport, error)
	Update(ctx context.Context, report *model.AnalyticsReport) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
	GetDueReports(ctx context.Context, now time.Time) ([]*model.AnalyticsReport, error)
	CreateDelivery(ctx context.Context, d *model.ReportDelivery) error
	ListDeliveries(ctx context.Context, reportID uuid.UUID, limit int) ([]*model.ReportDelivery, error)
}

type AnalyticsReportRepository struct {
	masterDB  *gorm.DB
	replicaDB *gorm.DB
}

func NewAnalyticsReportRepository(masterDB, replicaDB *gorm.DB) AnalyticsReportRepositoryI {
	return &AnalyticsReportRepository{masterDB: masterDB, replicaDB: replicaDB}
}

func (r *AnalyticsReportRepository) Create(ctx context.Context, report *model.AnalyticsReport) error {
	return r.masterDB.WithContext(ctx).Create(report).Error
}

func (r *AnalyticsReportRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (*model.AnalyticsReport, error) {
	var report model.AnalyticsReport
	err := r.replicaDB.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&report).Error
	return &report, err
}

func (r *AnalyticsReportRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*model.AnalyticsReport, error) {
	var reports []*model.AnalyticsReport
	err := r.replicaDB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&reports).Error
	return reports, err
}

func (r *AnalyticsReportRepository) Update(ctx context.Context, report *model.AnalyticsReport) error {
	return r.masterDB.WithContext(ctx).Save(report).Error
}

func (r *AnalyticsReportRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return r.masterDB.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.AnalyticsReport{}).Error
}

func (r *AnalyticsReportRepository) GetDueReports(ctx context.Context, now time.Time) ([]*model.AnalyticsReport, error) {
	var reports []*model.AnalyticsReport
	err := r.masterDB.WithContext(ctx).
		Where("is_active = true AND next_run_at <= ? AND deleted_at IS NULL", now).
		Find(&reports).Error
	return reports, err
}

func (r *AnalyticsReportRepository) CreateDelivery(ctx context.Context, d *model.ReportDelivery) error {
	return r.masterDB.WithContext(ctx).Create(d).Error
}

func (r *AnalyticsReportRepository) ListDeliveries(ctx context.Context, reportID uuid.UUID, limit int) ([]*model.ReportDelivery, error) {
	var deliveries []*model.ReportDelivery
	err := r.replicaDB.WithContext(ctx).
		Where("report_id = ?", reportID).
		Order("created_at DESC").
		Limit(limit).
		Find(&deliveries).Error
	return deliveries, err
}
