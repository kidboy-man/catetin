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

type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository implementation
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	model := r.domainToModel(user)

	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	if err := db.Create(model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrDuplicatePhoneNumber
		}
		return err
	}

	// Update domain entity with generated values
	user.ID = model.ID
	user.CreatedAt = model.CreatedAt
	user.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var model UserModel

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

func (r *userRepositoryImpl) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	var model UserModel

	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	if err := db.Where("phone_number = ?", phoneNumber).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return r.modelToDomain(&model), nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *domain.User) error {
	model := r.domainToModel(user)

	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	// Optimistic locking: check version
	result := db.Model(&UserModel{}).
		Where("id = ? AND version = ?", user.ID, user.Version-1).
		Updates(map[string]interface{}{
			"full_name":    model.FullName,
			"phone_number": model.PhoneNumber,
			"image":        model.Image,
			"version":      model.Version,
			"updated_at":   model.UpdatedAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrConflict
	}

	return nil
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	result := db.Delete(&UserModel{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *userRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var models []UserModel

	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	if err := db.Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	users := make([]*domain.User, len(models))
	for i, model := range models {
		users[i] = r.modelToDomain(&model)
	}

	return users, nil
}

// Helper methods for conversion between domain and model

func (r *userRepositoryImpl) domainToModel(user *domain.User) *UserModel {
	var deletedAt gorm.DeletedAt
	if user.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{
			Time:  *user.DeletedAt,
			Valid: true,
		}
	}

	return &UserModel{
		ID:          user.ID,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		Image:       user.Image,
		Version:     user.Version,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}

func (r *userRepositoryImpl) modelToDomain(model *UserModel) *domain.User {
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	return &domain.User{
		ID:          model.ID,
		FullName:    model.FullName,
		PhoneNumber: model.PhoneNumber,
		Image:       model.Image,
		Version:     model.Version,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}
