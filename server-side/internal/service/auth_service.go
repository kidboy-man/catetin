package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ingunawandra/catetin/internal/domain"
	"github.com/ingunawandra/catetin/internal/infrastructure/security"
	"github.com/ingunawandra/catetin/internal/repository"
	appErrors "github.com/ingunawandra/catetin/pkg/errors"
)

const EmailPasswordProviderName = "email-password"

// AuthService handles authentication business logic
type AuthService struct {
	userRepo         repository.UserRepository
	userAuthRepo     repository.UserAuthRepository
	authProviderRepo repository.AuthProviderRepository
	passwordHasher   *security.PasswordHasher
	jwtManager       *security.JWTManager
	txManager        repository.TransactionManager
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo repository.UserRepository,
	userAuthRepo repository.UserAuthRepository,
	authProviderRepo repository.AuthProviderRepository,
	passwordHasher *security.PasswordHasher,
	jwtManager *security.JWTManager,
	txManager repository.TransactionManager,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		userAuthRepo:     userAuthRepo,
		authProviderRepo: authProviderRepo,
		passwordHasher:   passwordHasher,
		jwtManager:       jwtManager,
		txManager:        txManager,
	}
}

// RegisterResponse represents the registration response
type RegisterResponse struct {
	User         *domain.User
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

// LoginResponse represents the login response
type LoginResponse struct {
	User         *domain.User
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

// Register registers a new user with email and password
func (s *AuthService) Register(ctx context.Context, fullName, email, password string) (*RegisterResponse, error) {
	// Get email-password auth provider
	provider, err := s.authProviderRepo.FindByName(ctx, EmailPasswordProviderName)
	if err != nil {
		return nil, appErrors.Wrap(err, appErrors.ErrCodeInternal, "Failed to find auth provider", 500)
	}
	if provider == nil {
		return nil, appErrors.New(appErrors.ErrCodeInternal, "Authentication provider not configured", 500)
	}

	// Check if email already exists
	existingAuth, err := s.userAuthRepo.FindByCredentialID(ctx, email, provider.ID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, appErrors.Wrap(err, appErrors.ErrCodeInternal, "Failed to check existing email", 500)
	}
	if existingAuth != nil {
		return nil, appErrors.ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := s.passwordHasher.Hash(password)
	if err != nil {
		return nil, appErrors.Wrap(err, appErrors.ErrCodeInternal, "Failed to hash password", 500)
	}

	var user *domain.User

	// Wrap user creation and auth creation in transaction
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// Create user
		user = domain.NewUser(fullName, email) // Use email as phone_number for now
		if err := s.userRepo.Create(txCtx, user); err != nil {
			return appErrors.Wrap(err, appErrors.ErrCodeInternal, "Failed to create user", 500)
		}

		// Create user auth
		userAuth := &repository.UserAuth{
			ID:               uuid.New(),
			UserID:           user.ID,
			AuthProviderID:   provider.ID,
			CredentialID:     email,
			CredentialSecret: hashedPassword,
		}
		if err := s.userAuthRepo.Create(txCtx, userAuth); err != nil {
			return appErrors.Wrap(err, appErrors.ErrCodeInternal, "Failed to create user auth", 500)
		}

		return nil // Commit transaction
	})

	if err != nil {
		return nil, err
	}

	// Generate tokens (outside transaction)
	accessToken, expiresIn, err := s.jwtManager.GenerateAccessToken(user.ID, email, fullName)
	if err != nil {
		return nil, appErrors.Wrap(err, appErrors.ErrCodeInternal, "Failed to generate access token", 500)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, appErrors.Wrap(err, appErrors.ErrCodeInternal, "Failed to generate refresh token", 500)
	}

	return &RegisterResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// Login authenticates a user with email and password
func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	// Get email-password auth provider
	provider, err := s.authProviderRepo.FindByName(ctx, EmailPasswordProviderName)
	if err != nil {
		return nil, appErrors.Wrap(err, appErrors.ErrCodeInternal, "Failed to find auth provider", 500)
	}
	if provider == nil {
		return nil, appErrors.New(appErrors.ErrCodeInternal, "Authentication provider not configured", 500)
	}

	// Find user auth by email
	userAuth, err := s.userAuthRepo.FindByCredentialID(ctx, email, provider.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, appErrors.ErrInvalidCredentials
		}
		return nil, appErrors.Wrap(err, appErrors.ErrCodeInternal, "Failed to find user auth", 500)
	}

	// Verify password
	if !s.passwordHasher.IsValidPassword(userAuth.CredentialSecret, password) {
		return nil, appErrors.ErrInvalidCredentials
	}

	// Get user details
	user, err := s.userRepo.FindByID(ctx, userAuth.UserID)
	if err != nil {
		return nil, appErrors.Wrap(err, appErrors.ErrCodeInternal, "Failed to find user", 500)
	}

	// Generate tokens
	accessToken, expiresIn, err := s.jwtManager.GenerateAccessToken(user.ID, email, user.FullName)
	if err != nil {
		return nil, appErrors.Wrap(err, appErrors.ErrCodeInternal, "Failed to generate access token", 500)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, appErrors.Wrap(err, appErrors.ErrCodeInternal, "Failed to generate refresh token", 500)
	}

	return &LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// EnsureEmailPasswordProvider ensures the email-password auth provider exists
func (s *AuthService) EnsureEmailPasswordProvider(ctx context.Context) error {
	provider, err := s.authProviderRepo.FindByName(ctx, EmailPasswordProviderName)
	if err != nil {
		return fmt.Errorf("failed to check auth provider: %w", err)
	}

	if provider == nil {
		// Create the email-password provider
		name := EmailPasswordProviderName
		provider = &repository.AuthProvider{
			ID:          uuid.New(),
			DisplayName: "Email & Password",
			Name:        &name,
		}
		if err := s.authProviderRepo.Create(ctx, provider); err != nil {
			return fmt.Errorf("failed to create auth provider: %w", err)
		}
	}

	return nil
}
