package user

import (
	"context"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"github.com/deni12345/dae-services/services/dae-core/internal/port"
	"github.com/deni12345/dae-services/libs/apperror"
)

// GetUser retrieves a user by ID
func (u *usecase) GetUser(ctx context.Context, id string) (*domain.User, error) {
	ctx, span := tracer.Start(ctx, "UserUC.GetUser")
	defer span.End()

	if id == "" {
		err := apperror.InvalidInput("user_id is required")
		span.RecordError(err)
		return nil, err
	}

	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return user, nil
}

// ListUsers retrieves a paginated list of users
func (u *usecase) ListUsers(ctx context.Context, req *ListUsersReq) (*ListUsersResp, error) {
	ctx, span := tracer.Start(ctx, "UserUC.ListUsers")
	defer span.End()

	// Set default limit
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// Fetch one extra to determine if there are more results
	users, err := u.userRepo.List(ctx, port.ListUserQuery{
		Limit:           req.PageSize + 1,
		Cursor:          req.Cursor,
		IncludeDisabled: req.IncludeDisabled,
	})
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Determine if there are more results
	var nextCursor string
	if int32(len(users)) > req.PageSize {
		// Trim to requested page size
		users = users[:req.PageSize]
		// Use last item's ID as next cursor
		nextCursor = users[len(users)-1].ID
	}

	return &ListUsersResp{
		Users:      users,
		NextCursor: nextCursor,
	}, nil
}
