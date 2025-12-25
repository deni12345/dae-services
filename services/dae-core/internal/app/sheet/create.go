package sheet

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"github.com/deni12345/dae-services/services/dae-core/internal/grpc/interceptor"
	"github.com/deni12345/dae-services/libs/apperror"
	"github.com/google/uuid"
)

const (
	idempotencyTTL = 24 * time.Hour
)

// CreateSheet creates a new sheet with idempotency protection
func (u *usecase) CreateSheet(ctx context.Context, req *CreateSheetReq) (*domain.Sheet, error) {
	ctx, span := tracer.Start(ctx, "SheetUC.CreateSheet")
	defer span.End()

	if err := validateCreateRequest(req); err != nil {
		span.RecordError(err)
		return nil, err
	}

	idemKey := interceptor.GetOrCreateIdempotencyKeyWithHash(ctx, req.HostUserID)

	result, err := u.idemStore.Do(ctx, idemKey, idempotencyTTL, func(ctx context.Context) ([]byte, error) {
		sheet, err := u.createSheetInternal(ctx, req)
		if err != nil {
			return nil, err
		}

		return json.Marshal(sheet)
	})

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	var sheet domain.Sheet
	if err := json.Unmarshal(result, &sheet); err != nil {
		span.RecordError(err)
		return nil, apperror.Internal(fmt.Sprintf("unmarshal sheet: %v", err))
	}

	return &sheet, nil
}

// validateCreateRequest validates the create sheet request
func validateCreateRequest(req *CreateSheetReq) error {
	if req.HostUserID == "" {
		return apperror.InvalidInput("host_user_id is required")
	}
	if req.Name == "" {
		return apperror.InvalidInput("name is required")
	}
	return nil
}

// createSheetInternal is the core sheet creation logic
func (u *usecase) createSheetInternal(ctx context.Context, req *CreateSheetReq) (*domain.Sheet, error) {
	now := time.Now().UTC()

	// Validate menu items BEFORE conversion
	if len(req.MenuItems) > 0 {
		if err := validateMenuItems(req.MenuItems); err != nil {
			return nil, err
		}
	}

	// Initialize member list with host as first member
	memberIDs := []string{req.HostUserID}
	if req.MemberIDs != nil {
		// Add additional members, avoiding duplicates
		seen := map[string]bool{req.HostUserID: true}
		for _, mid := range req.MemberIDs {
			if !seen[mid] {
				memberIDs = append(memberIDs, mid)
				seen[mid] = true
			}
		}
	}

	sheet := &domain.Sheet{
		ID:          fmt.Sprintf("%s-%s", req.Name, uuid.New().String()),
		Name:        req.Name,
		Description: req.Description,
		HostUserID:  req.HostUserID,
		Status:      domain.Status_OPEN, // Default to open
		DeliveryFee: *req.DeliveryFee,
		Discount:    req.Discount,
		MemberIDs:   memberIDs,

		CreatedAt: now,
		UpdatedAt: now,
	}

	createdSheet, err := u.sheetRepo.Create(ctx, sheet)
	if err != nil {
		return nil, err
	}

	// Add members to subcollection (dual-write pattern)
	for _, memberID := range memberIDs {
		if err := u.sheetRepo.AddMember(ctx, sheet.ID, memberID); err != nil {
			return nil, err
		}
	}

	// Attach menu items if provided (already validated)
	if len(req.MenuItems) > 0 {
		menuItems := convertMenuItemsToDomain(req.MenuItems, now.Unix())
		if err := u.sheetRepo.AttachMenuItems(ctx, sheet.ID, menuItems); err != nil {
			return nil, err
		}
	}

	return createdSheet, nil
}

// validateMenuItems validates all menu items and their nested structures
// Returns early on first validation error
func validateMenuItems(items []MenuItemReq) error {
	if len(items) == 0 {
		return nil
	}

	seen := make(map[string]bool, len(items))
	for _, item := range items {
		// Item name required
		if item.Name == "" {
			return ErrMenuItemNameRequired
		}

		// Duplicate detection
		if seen[item.Name] {
			return ErrDuplicateMenuItemName
		}
		seen[item.Name] = true

		// Price & currency validation
		if item.Price < 0 {
			return ErrMenuItemInvalidPrice
		}
		if item.Currency == "" {
			return ErrMenuItemInvalidCurrency
		}

		// Validate option groups
		if err := validateOptionGroups(item.OptionGroups); err != nil {
			return err
		}
	}
	return nil
}

func validateOptionGroups(groups []MenuOptionGroupReq) error {
	if len(groups) == 0 {
		return nil
	}

	seen := make(map[string]bool, len(groups))
	for _, grp := range groups {
		if grp.Name == "" {
			return ErrOptionGroupNameRequired
		}
		if seen[grp.Name] {
			return ErrDuplicateOptionGroupName
		}
		seen[grp.Name] = true

		// MaxSelect validation for multi-select groups
		if grp.MultiSelect && grp.MaxSelect <= 0 {
			return ErrOptionGroupInvalidMaxSelect
		}

		// Validate options
		if err := validateOptions(grp.Options); err != nil {
			return err
		}
	}
	return nil
}

func validateOptions(options []MenuOptionReq) error {
	if len(options) == 0 {
		return nil
	}

	seen := make(map[string]bool, len(options))
	for _, opt := range options {
		if opt.Name == "" {
			return ErrOptionNameRequired
		}
		if seen[opt.Name] {
			return ErrDuplicateOptionName
		}
		seen[opt.Name] = true

		if opt.Price < 0 {
			return ErrOptionInvalidPrice
		}
	}
	return nil
}

// convertMenuItemsToDomain converts validated request DTOs to domain entities
// Assumes items are already validated by validateMenuItems
func convertMenuItemsToDomain(reqItems []MenuItemReq, timestamp int64) []*domain.MenuItem {
	domainItems := make([]*domain.MenuItem, len(reqItems))

	for i, req := range reqItems {
		domainItems[i] = &domain.MenuItem{
			ID:           fmt.Sprintf("%s-%s", req.Name, uuid.New().String()),
			Name:         req.Name,
			Active:       req.Active,
			Price:        req.Price,
			Currency:     req.Currency,
			OptionGroups: convertOptionGroups(req.OptionGroups),
			UpdatedAt:    timestamp,
		}
	}
	return domainItems
}

// convertOptionGroups is extracted to reduce nesting
func convertOptionGroups(reqGroups []MenuOptionGroupReq) map[string]domain.OptionGroup {
	if len(reqGroups) == 0 {
		return nil // Return nil instead of empty map (saves memory)
	}

	optionGroups := make(map[string]domain.OptionGroup, len(reqGroups))
	for _, grpReq := range reqGroups {
		id := fmt.Sprintf("%s-%s", grpReq.Name, uuid.New().String())

		// Convert MultiSelect bool to domain.OptionGroupType
		var groupType domain.OptionGroupType
		if grpReq.MultiSelect {
			groupType = domain.GroupMulti
		} else {
			groupType = domain.GroupSingle
		}

		optionGroups[id] = domain.OptionGroup{
			ID:        id,
			Name:      grpReq.Name,
			Type:      groupType,
			Required:  grpReq.Required,
			MaxSelect: int(grpReq.MaxSelect),
			Options:   convertOptions(grpReq.Options),
		}
	}
	return optionGroups
}

// convertOptions is extracted for clarity
func convertOptions(reqOptions []MenuOptionReq) map[string]domain.Option {
	if len(reqOptions) == 0 {
		return nil
	}

	options := make(map[string]domain.Option, len(reqOptions))
	for _, optReq := range reqOptions {
		id := fmt.Sprintf("%s-%s", optReq.Name, uuid.New().String())
		options[id] = domain.Option{
			ID:     id,
			Name:   optReq.Name,
			Price:  optReq.Price,
			Per:    domain.PerUnit, // Default to per-unit pricing
			Active: optReq.Active,
		}
	}
	return options
}
