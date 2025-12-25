package user

import (
	"context"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"github.com/deni12345/dae-services/services/dae-core/internal/port"
	"github.com/deni12345/dae-services/libs/utils"
)

func (uc *usecase) CreateUser(ctx context.Context, req *CreateUserReq) (*domain.User, error) {
	ctx, span := tracer.Start(ctx, "UserUC.CreateUser")
	defer span.End()

	// Normalize email
	normalizedEmail := utils.NormalizeString(req.Email)

	// Check if email is already taken (pre-check for better error messages)
	unique, err := uc.userRepo.CheckEmailUnique(ctx, normalizedEmail)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	if !unique {
		err := ErrEmailAlreadyExists
		span.RecordError(err)
		return nil, err
	}

	// Check if identity is already linked
	unique, err = uc.userRepo.CheckIdentityUnique(ctx, req.Provider, req.Subject)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	if !unique {
		err := ErrIdentityAlreadyLinked
		span.RecordError(err)
		return nil, err
	}

	// Create user
	createReq := port.CreateUserRequest{
		Email:       req.Email,
		Name:        req.Name,
		DisplayName: req.DisplayName,
		PhotoURL:    req.PhotoURL,
		Phone:       req.Phone,
		Provider:    req.Provider,
		Subject:     req.Subject,
		Password:    req.Password,
	}

	user, err := uc.userRepo.Create(ctx, createReq)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return user, nil
}
