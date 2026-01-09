package postgresql

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ingunawandra/catetin/internal/domain"
	"github.com/ingunawandra/catetin/internal/repository"
	"gorm.io/gorm"
)

type moneyFlowRepositoryImpl struct {
	db *gorm.DB
}

// NewMoneyFlowRepository creates a new money flow repository implementation
func NewMoneyFlowRepository(db *gorm.DB) repository.MoneyFlowRepository {
	return &moneyFlowRepositoryImpl{db: db}
}

func (r *moneyFlowRepositoryImpl) Create(ctx context.Context, moneyFlow *domain.MoneyFlow) error {
	model := r.domainToModel(moneyFlow)

	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	if err := db.Create(model).Error; err != nil {
		return err
	}

	// Update domain entity with generated values
	moneyFlow.ID = model.ID
	moneyFlow.CreatedAt = model.CreatedAt
	moneyFlow.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *moneyFlowRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.MoneyFlow, error) {
	var model MoneyFlowModel

	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	if err := db.Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return r.modelToDomain(&model), nil
}

func (r *moneyFlowRepositoryImpl) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.MoneyFlow, error) {
	var models []MoneyFlowModel

	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	if err := db.Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	moneyFlows := make([]*domain.MoneyFlow, len(models))
	for i, model := range models {
		moneyFlows[i] = r.modelToDomain(&model)
	}

	return moneyFlows, nil
}

func (r *moneyFlowRepositoryImpl) FindByUserIDAndDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*domain.MoneyFlow, error) {
	var models []MoneyFlowModel

	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	if err := db.Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	moneyFlows := make([]*domain.MoneyFlow, len(models))
	for i, model := range models {
		moneyFlows[i] = r.modelToDomain(&model)
	}

	return moneyFlows, nil
}

func (r *moneyFlowRepositoryImpl) Update(ctx context.Context, moneyFlow *domain.MoneyFlow) error {
	model := r.domainToModel(moneyFlow)

	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	// Optimistic locking: check version
	result := db.Model(&MoneyFlowModel{}).
		Where("id = ? AND version = ?", moneyFlow.ID, moneyFlow.Version-1).
		Updates(map[string]any{
			"category":    model.Category,
			"amount":      model.Amount,
			"currency":    model.Currency,
			"description": model.Description,
			"tags":        model.Tags,
			"version":     model.Version,
			"updated_at":  model.UpdatedAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrConflict
	}

	return nil
}

func (r *moneyFlowRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	result := db.Delete(&MoneyFlowModel{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *moneyFlowRepositoryImpl) GetTotalByUserID(ctx context.Context, userID uuid.UUID) (float64, error) {
	var total float64

	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	if err := db.Model(&MoneyFlowModel{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *moneyFlowRepositoryImpl) GetTotalByUserIDAndCategory(ctx context.Context, userID uuid.UUID, category string) (float64, error) {
	var total float64

	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	if err := db.Model(&MoneyFlowModel{}).
		Where("user_id = ? AND category = ?", userID, category).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// Helper methods for conversion between domain and model

func (r *moneyFlowRepositoryImpl) domainToModel(moneyFlow *domain.MoneyFlow) *MoneyFlowModel {
	var deletedAt gorm.DeletedAt
	if moneyFlow.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{
			Time:  *moneyFlow.DeletedAt,
			Valid: true,
		}
	}

	tags := JSONB(moneyFlow.Tags)
	if tags == nil {
		tags = JSONB([]string{})
	}

	return &MoneyFlowModel{
		ID:          moneyFlow.ID,
		UserID:      moneyFlow.UserID,
		Category:    moneyFlow.Category,
		Amount:      moneyFlow.Amount,
		Currency:    moneyFlow.Currency,
		Description: moneyFlow.Description,
		Tags:        tags,
		Version:     moneyFlow.Version,
		CreatedAt:   moneyFlow.CreatedAt,
		UpdatedAt:   moneyFlow.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}

func (r *moneyFlowRepositoryImpl) modelToDomain(model *MoneyFlowModel) *domain.MoneyFlow {
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	tags := []string(model.Tags)
	if tags == nil {
		tags = []string{}
	}

	return &domain.MoneyFlow{
		ID:          model.ID,
		UserID:      model.UserID,
		Category:    model.Category,
		Amount:      model.Amount,
		Currency:    model.Currency,
		Description: model.Description,
		Tags:        tags,
		Version:     model.Version,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}
