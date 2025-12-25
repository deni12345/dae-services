package user

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"github.com/deni12345/dae-services/services/dae-core/internal/port"
	"github.com/deni12345/dae-services/libs/utils"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Create creates a new user and enforces uniqueness via a transaction.
func (r *userRepo) Create(ctx context.Context, req port.CreateUserRequest) (*domain.User, error) {
	ctx, span := tracer.Start(ctx, "UserRepo.Create")
	defer span.End()

	normalizedEmail := utils.NormalizeString(req.Email)

	// Pre-validate before transaction
	if normalizedEmail == "" {
		err := fmt.Errorf("email is required")
		span.RecordError(err)
		return nil, err
	}
	if req.Name == "" {
		err := fmt.Errorf("name is required")
		span.RecordError(err)
		return nil, err
	}
	if req.Provider == "" {
		err := fmt.Errorf("provider is required")
		span.RecordError(err)
		return nil, err
	}
	if req.Subject == "" {
		err := fmt.Errorf("subject is required")
		span.RecordError(err)
		return nil, err
	}

	// Hash password if local provider
	var passwordHash, passwordAlgo *string
	var passwordUpdatedAt *time.Time
	if req.Provider == domain.IdentityProviderLocal && req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("hash password: %w", err)
		}
		hashStr := string(hash)
		algoStr := "bcrypt"
		now := time.Now().UTC()
		passwordHash = &hashStr
		passwordAlgo = &algoStr
		passwordUpdatedAt = &now
	}

	now := time.Now().UTC()
	displayName := req.DisplayName
	if displayName == "" {
		displayName = req.Name
	}

	var createdUser *domain.User

	err := r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// Ensure uniqueness and prepare to create user and identity entries
		uniqueEmailRef := r.client.Collection("unique_emails").Doc(normalizedEmail)
		uniqueEmailSnap, err := tx.Get(uniqueEmailRef)
		if err != nil && status.Code(err) != codes.NotFound {
			return fmt.Errorf("check email uniqueness: %w", err)
		}
		if uniqueEmailSnap.Exists() {
			return ErrAlreadyExists
		}

		identityKey := fmt.Sprintf("%s:%s", req.Provider, req.Subject)
		uniqueIdentityRef := r.client.Collection("unique_identities").Doc(identityKey)
		uniqueIdentitySnap, err := tx.Get(uniqueIdentityRef)
		if err != nil && status.Code(err) != codes.NotFound {
			return fmt.Errorf("check identity uniqueness: %w", err)
		}
		if uniqueIdentitySnap.Exists() {
			return ErrAlreadyExists
		}

		// Create user document
		userRef := r.collection.NewDoc()
		userID := userRef.ID

		user := &domain.User{
			ID:                userID,
			Email:             req.Email,
			EmailNormalized:   normalizedEmail,
			EmailVerified:     req.Provider != domain.IdentityProviderLocal, // Auto-verify for OAuth providers
			Name:              req.Name,
			DisplayName:       displayName,
			PhotoURL:          req.PhotoURL,
			Phone:             req.Phone,
			Roles:             []domain.Role{domain.RoleUser},
			Status:            domain.UserStatusActive,
			CreatedAt:         now,
			UpdatedAt:         now,
			PasswordHash:      passwordHash,
			PasswordAlgo:      passwordAlgo,
			PasswordUpdatedAt: passwordUpdatedAt,
		}

		if err := tx.Set(userRef, user); err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		// 4. Create identity subcollection document
		identityRef := userRef.Collection("identities").Doc(string(req.Provider))
		identity := &domain.UserIdentity{
			UserID:        userID,
			Provider:      req.Provider,
			Subject:       req.Subject,
			EmailAtSignup: req.Email,
			LinkedAt:      now,
		}
		if err := tx.Set(identityRef, identity); err != nil {
			return fmt.Errorf("create identity: %w", err)
		}

		// 5. Create unique_emails entry
		uniqueEmail := &domain.UniqueEmail{
			UserID:    userID,
			CreatedAt: now,
		}
		if err := tx.Set(uniqueEmailRef, uniqueEmail); err != nil {
			return fmt.Errorf("create unique email: %w", err)
		}

		// 6. Create unique_identities entry
		uniqueIdentity := &domain.UniqueIdentity{
			UserID:    userID,
			CreatedAt: now,
		}
		if err := tx.Set(uniqueIdentityRef, uniqueIdentity); err != nil {
			return fmt.Errorf("create unique identity: %w", err)
		}

		createdUser = user
		return nil
	})

	if err != nil {
		span.RecordError(err)
		if err == ErrAlreadyExists {
			return nil, err
		}
		return nil, fmt.Errorf("create user transaction: %w", err)
	}

	return createdUser, nil
}

// CreateIdentity adds a new identity provider to an existing user
func (r *userRepo) CreateIdentity(ctx context.Context, userID string, identity *domain.UserIdentity) error {
	ctx, span := tracer.Start(ctx, "UserRepo.CreateIdentity")
	defer span.End()

	err := r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// Check if identity already exists
		identityKey := fmt.Sprintf("%s:%s", identity.Provider, identity.Subject)
		uniqueIdentityRef := r.client.Collection("unique_identities").Doc(identityKey)
		uniqueIdentitySnap, err := tx.Get(uniqueIdentityRef)
		if err != nil && status.Code(err) != codes.NotFound {
			return fmt.Errorf("check identity uniqueness: %w", err)
		}
		if uniqueIdentitySnap.Exists() {
			return ErrAlreadyExists
		}

		// Create identity subcollection document
		userRef := r.collection.Doc(userID)
		identityRef := userRef.Collection("identities").Doc(string(identity.Provider))
		if err := tx.Set(identityRef, identity); err != nil {
			return fmt.Errorf("create identity: %w", err)
		}

		// Create unique_identities entry
		now := time.Now().UTC()
		uniqueIdentity := &domain.UniqueIdentity{
			UserID:    userID,
			CreatedAt: now,
		}
		if err := tx.Set(uniqueIdentityRef, uniqueIdentity); err != nil {
			return fmt.Errorf("create unique identity: %w", err)
		}

		return nil
	})

	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("create identity transaction: %w", err)
	}

	return nil
}

// GetIdentityByProvider retrieves an identity for a user by provider
func (r *userRepo) GetIdentityByProvider(ctx context.Context, userID string, provider domain.IdentityProvider) (*domain.UserIdentity, error) {
	ctx, span := tracer.Start(ctx, "UserRepo.GetIdentityByProvider")
	defer span.End()

	identityRef := r.collection.Doc(userID).Collection("identities").Doc(string(provider))
	snap, err := identityRef.Get(ctx)
	if err != nil {
		span.RecordError(err)
		if status.Code(err) == codes.NotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get identity: %w", err)
	}

	var identity domain.UserIdentity
	if err := snap.DataTo(&identity); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("unmarshal identity: %w", err)
	}
	identity.ID = snap.Ref.ID

	return &identity, nil
}

// CheckEmailUnique checks if an email is available
func (r *userRepo) CheckEmailUnique(ctx context.Context, email string) (bool, error) {
	ctx, span := tracer.Start(ctx, "UserRepo.CheckEmailUnique")
	defer span.End()

	normalizedEmail := utils.NormalizeString(email)
	uniqueEmailRef := r.client.Collection("unique_emails").Doc(normalizedEmail)
	snap, err := uniqueEmailRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return true, nil // Email is unique
		}
		span.RecordError(err)
		return false, fmt.Errorf("check email uniqueness: %w", err)
	}

	return !snap.Exists(), nil
}

// CheckIdentityUnique checks if a provider identity is available
func (r *userRepo) CheckIdentityUnique(ctx context.Context, provider domain.IdentityProvider, subject string) (bool, error) {
	ctx, span := tracer.Start(ctx, "UserRepo.CheckIdentityUnique")
	defer span.End()

	identityKey := fmt.Sprintf("%s:%s", provider, subject)
	uniqueIdentityRef := r.client.Collection("unique_identities").Doc(identityKey)
	snap, err := uniqueIdentityRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return true, nil // Identity is unique
		}
		span.RecordError(err)
		return false, fmt.Errorf("check identity uniqueness: %w", err)
	}

	return !snap.Exists(), nil
}
