package postgresql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ingunawandra/catetin/internal/domain"
	"github.com/ingunawandra/catetin/internal/repository"
	"gorm.io/gorm"
)

type userAuthRepositoryImpl struct {
	db repository.DB
}

// NewUserAuthRepository creates a new user auth repository implementation
func NewUserAuthRepository(db repository.DB) repository.UserAuthRepository {
	return &userAuthRepositoryImpl{db: db}
}

func (r *userAuthRepositoryImpl) Create(ctx context.Context, userAuth *repository.UserAuth) error {
	model := r.domainToModel(userAuth)

	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	res := db.Create(model)
	if err := res.Error(); err != nil {
		return err
	}

	userAuth.ID = model.ID
	return nil
}

func (r *userAuthRepositoryImpl) FindByCredentialID(ctx context.Context, credentialID string, authProviderID uuid.UUID) (*repository.UserAuth, error) {
	var model UserAuthModel

	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	res := db.Where("credential_id = ? AND auth_provider_id = ?", credentialID, authProviderID).
		First(&model)
	if err := res.Error(); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return r.modelToDomain(&model), nil
}

func (r *userAuthRepositoryImpl) FindByUserIDAndProvider(ctx context.Context, userID, authProviderID uuid.UUID) (*repository.UserAuth, error) {
	var model UserAuthModel

	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	res := db.Where("user_id = ? AND auth_provider_id = ?", userID, authProviderID).
		First(&model)
	if err := res.Error(); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return r.modelToDomain(&model), nil
}

func (r *userAuthRepositoryImpl) Update(ctx context.Context, userAuth *repository.UserAuth) error {
	model := r.domainToModel(userAuth)

	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	result := db.Model(&UserAuthModel{}).
		Where("id = ?", userAuth.ID).
		Updates(map[string]interface{}{
			"credential_id":      model.CredentialID,
			"credential_secret":  model.CredentialSecret,
			"credential_refresh": model.CredentialRefresh,
			"updated_at":         model.UpdatedAt,
		})

	if err := result.Error(); err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *userAuthRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	// Use GetDB to support transactions
	db := GetDB(ctx, r.db)

	result := db.Delete(&UserAuthModel{}, "id = ?", id)

	if err := result.Error(); err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Helper methods for conversion

func (r *userAuthRepositoryImpl) domainToModel(userAuth *repository.UserAuth) *UserAuthModel {
	return &UserAuthModel{
		ID:                userAuth.ID,
		UserID:            userAuth.UserID,
		AuthProviderID:    userAuth.AuthProviderID,
		CredentialID:      userAuth.CredentialID,
		CredentialSecret:  userAuth.CredentialSecret,
		CredentialRefresh: userAuth.CredentialRefresh,
	}
}

func (r *userAuthRepositoryImpl) modelToDomain(model *UserAuthModel) *repository.UserAuth {
	return &repository.UserAuth{
		ID:                model.ID,
		UserID:            model.UserID,
		AuthProviderID:    model.AuthProviderID,
		CredentialID:      model.CredentialID,
		CredentialSecret:  model.CredentialSecret,
		CredentialRefresh: model.CredentialRefresh,
	}
}
