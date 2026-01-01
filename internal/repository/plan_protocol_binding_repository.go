package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// PlanProtocolBinding links a plan to protocol bindings.
type PlanProtocolBinding struct {
	PlanID    uint64    `gorm:"primaryKey;autoIncrement:false;index"`
	BindingID uint64    `gorm:"primaryKey;autoIncrement:false;index"`
	CreatedAt time.Time
}

// TableName binds the mapping table name.
func (PlanProtocolBinding) TableName() string { return "plan_protocol_bindings" }

// PlanProtocolBindingRepository manages plan to binding mappings.
type PlanProtocolBindingRepository interface {
	ListBindingIDs(ctx context.Context, planID uint64) ([]uint64, error)
	ListBindingsByPlanIDs(ctx context.Context, planIDs []uint64) (map[uint64][]uint64, error)
	Replace(ctx context.Context, planID uint64, bindingIDs []uint64) error
}

type planProtocolBindingRepository struct {
	db *gorm.DB
}

// NewPlanProtocolBindingRepository constructs a mapping repository.
func NewPlanProtocolBindingRepository(db *gorm.DB) (PlanProtocolBindingRepository, error) {
	if db == nil {
		return nil, errors.New("repository: database connection is required")
	}
	return &planProtocolBindingRepository{db: db}, nil
}

func (r *planProtocolBindingRepository) ListBindingIDs(ctx context.Context, planID uint64) ([]uint64, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if planID == 0 {
		return nil, ErrInvalidArgument
	}

	var rows []PlanProtocolBinding
	if err := r.db.WithContext(ctx).
		Where("plan_id = ?", planID).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]uint64, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.BindingID)
	}
	return result, nil
}

func (r *planProtocolBindingRepository) ListBindingsByPlanIDs(ctx context.Context, planIDs []uint64) (map[uint64][]uint64, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(planIDs) == 0 {
		return map[uint64][]uint64{}, nil
	}

	var rows []PlanProtocolBinding
	if err := r.db.WithContext(ctx).
		Where("plan_id IN ?", planIDs).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make(map[uint64][]uint64, len(planIDs))
	for _, row := range rows {
		result[row.PlanID] = append(result[row.PlanID], row.BindingID)
	}
	return result, nil
}

func (r *planProtocolBindingRepository) Replace(ctx context.Context, planID uint64, bindingIDs []uint64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if planID == 0 {
		return ErrInvalidArgument
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("plan_id = ?", planID).Delete(&PlanProtocolBinding{}).Error; err != nil {
			return err
		}
		if len(bindingIDs) == 0 {
			return nil
		}

		now := time.Now().UTC()
		rows := make([]PlanProtocolBinding, 0, len(bindingIDs))
		for _, id := range bindingIDs {
			if id == 0 {
				return ErrInvalidArgument
			}
			rows = append(rows, PlanProtocolBinding{
				PlanID:    planID,
				BindingID: id,
				CreatedAt: now,
			})
		}

		return tx.Create(&rows).Error
	})
}
