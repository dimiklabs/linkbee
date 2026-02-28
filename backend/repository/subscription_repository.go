package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/shafikshaon/linkbee/model"
)

type SubscriptionRepositoryI interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Subscription, error)
	GetByLemonSqueezySubID(ctx context.Context, lsSubID string) (*model.Subscription, error)
	Upsert(ctx context.Context, sub *model.Subscription) error
}

type SubscriptionRepository struct {
	masterDB  *gorm.DB
	replicaDB *gorm.DB
}

func NewSubscriptionRepository(masterDB, replicaDB *gorm.DB) SubscriptionRepositoryI {
	return &SubscriptionRepository{masterDB: masterDB, replicaDB: replicaDB}
}

func (r *SubscriptionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Subscription, error) {
	var sub model.Subscription
	if err := r.replicaDB.WithContext(ctx).Where("user_id = ?", userID).First(&sub).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *SubscriptionRepository) GetByLemonSqueezySubID(ctx context.Context, lsSubID string) (*model.Subscription, error) {
	var sub model.Subscription
	if err := r.replicaDB.WithContext(ctx).Where("lemon_squeezy_sub_id = ?", lsSubID).First(&sub).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

// Upsert inserts or updates based on user_id uniqueness.
func (r *SubscriptionRepository) Upsert(ctx context.Context, sub *model.Subscription) error {
	return r.masterDB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"plan_id", "status",
				"lemon_squeezy_sub_id", "lemon_squeezy_order_id",
				"lemon_squeezy_customer_id", "lemon_squeezy_variant_id",
				"current_period_end", "cancelled_at", "updated_at",
			}),
		}).
		Create(sub).Error
}
