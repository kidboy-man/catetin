package postgresql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ingunawandra/catetin/internal/repository"
	"gorm.io/gorm"
)

type authProviderRepositoryImpl struct {
	db *gorm.DB
}

// NewAuthProviderRepository creates a new auth provider repository implementation
func NewAuthProviderRepository(db *gorm.DB) repository.AuthProviderRepository {
	return &authProviderRepositoryImpl{db: db}
}

func (r *authProviderRepositoryImpl) FindByName(ctx context.Context, name string) (*repository.AuthProvider, error) {
	var model AuthProviderModel

	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	if err := db.Where("name = ?", name).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found (not an error for this case)
		}
		return nil, err
	}

	return r.modelToDomain(&model), nil
}

func (r *authProviderRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*repository.AuthProvider, error) {
	var model AuthProviderModel

	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	if err := db.Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToDomain(&model), nil
}

func (r *authProviderRepositoryImpl) Create(ctx context.Context, provider *repository.AuthProvider) error {
	model := r.domainToModel(provider)

	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	if err := db.Create(model).Error; err != nil {
		return err
	}

	provider.ID = model.ID
	return nil
}

// Helper methods for conversion

func (r *authProviderRepositoryImpl) domainToModel(provider *repository.AuthProvider) *AuthProviderModel {
	return &AuthProviderModel{
		ID:           provider.ID,
		DisplayName:  provider.DisplayName,
		Name:         provider.Name,
		Image:        provider.Image,
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
	}
}

func (r *authProviderRepositoryImpl) modelToDomain(model *AuthProviderModel) *repository.AuthProvider {
	return &repository.AuthProvider{
		ID:           model.ID,
		DisplayName:  model.DisplayName,
		Name:         model.Name,
		Image:        model.Image,
		ClientID:     model.ClientID,
		ClientSecret: model.ClientSecret,
	}
}
