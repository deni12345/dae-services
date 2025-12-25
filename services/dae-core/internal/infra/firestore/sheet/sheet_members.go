package sheet

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"google.golang.org/api/iterator"
)

const (
	MemberRoleHost   = "host"
	MemberRoleMember = "member"
)

// AddMember adds a member (denormalized + subcollection)
func (r *sheetRepo) AddMember(ctx context.Context, sheetID, userID string) error {
	ctx, span := tracer.Start(ctx, "SheetRepo.AddMember")
	defer span.End()

	err := r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		sheetRef := r.collection.Doc(sheetID)
		snap, err := tx.Get(sheetRef)
		if err != nil {
			return fmt.Errorf("get sheet: %w", err)
		}

		var sheet domain.Sheet
		if err := snap.DataTo(&sheet); err != nil {
			return fmt.Errorf("unmarshal sheet: %w", err)
		}

		for _, id := range sheet.MemberIDs {
			if id == userID {
				return nil
			}
		}

		sheet.MemberIDs = append(sheet.MemberIDs, userID)
		now := time.Now().UTC()

		updates := []firestore.Update{
			{Path: "member_ids", Value: sheet.MemberIDs},
			{Path: "updated_at", Value: now},
		}

		if err := tx.Update(sheetRef, updates); err != nil {
			return fmt.Errorf("update sheet members: %w", err)
		}

		memberRef := sheetRef.Collection("members").Doc(userID)
		memberData := map[string]interface{}{
			"user_id":   userID,
			"role":      MemberRoleMember,
			"joined_at": now,
		}

		return tx.Set(memberRef, memberData)
	})

	if err != nil {
		span.RecordError(err)
	}
	return err
}

// RemoveMember removes a user from a sheet's member list
func (r *sheetRepo) RemoveMember(ctx context.Context, sheetID, userID string) error {
	ctx, span := tracer.Start(ctx, "SheetRepo.RemoveMember")
	defer span.End()

	err := r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		sheetRef := r.collection.Doc(sheetID)
		snap, err := tx.Get(sheetRef)
		if err != nil {
			return fmt.Errorf("get sheet: %w", err)
		}

		var sheet domain.Sheet
		if err := snap.DataTo(&sheet); err != nil {
			return fmt.Errorf("unmarshal sheet: %w", err)
		}

		// Remove from denormalized array
		newMembers := make([]string, 0, len(sheet.MemberIDs))
		found := false
		for _, id := range sheet.MemberIDs {
			if id != userID {
				newMembers = append(newMembers, id)
			} else {
				found = true
			}
		}

		if !found {
			return nil // not a member, idempotent
		}

		now := time.Now().UTC()
		updates := []firestore.Update{
			{Path: "member_ids", Value: newMembers},
			{Path: "updated_at", Value: now},
		}

		if err := tx.Update(sheetRef, updates); err != nil {
			return fmt.Errorf("update sheet members: %w", err)
		}

		// Remove from subcollection
		memberRef := sheetRef.Collection("members").Doc(userID)
		return tx.Delete(memberRef)
	})

	if err != nil {
		span.RecordError(err)
	}
	return err
}

// ListMemberIDs returns all member IDs for a sheet
func (r *sheetRepo) ListMemberIDs(ctx context.Context, sheetID string) ([]string, error) {
	ctx, span := tracer.Start(ctx, "SheetRepo.ListMemberIDs")
	defer span.End()

	snap, err := r.collection.Doc(sheetID).Get(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("get sheet: %w", err)
	}

	var sheet domain.Sheet
	if err := snap.DataTo(&sheet); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("unmarshal sheet: %w", err)
	}

	// Return denormalized list for fast read
	if sheet.MemberIDs == nil {
		return []string{}, nil
	}

	return sheet.MemberIDs, nil
}

// syncMemberToSubcollection is a helper to ensure subcollection is in sync
// Use this during migration or repair operations
func (r *sheetRepo) syncMemberToSubcollection(ctx context.Context, sheetID, userID string, role string, joinedAt time.Time) error {
	memberRef := r.collection.Doc(sheetID).Collection("members").Doc(userID)
	_, err := memberRef.Set(ctx, map[string]interface{}{
		"user_id":   userID,
		"role":      role,
		"joined_at": joinedAt,
	})
	return err
}

// listMembersFromSubcollection returns members from subcollection (for verification/debug)
func (r *sheetRepo) listMembersFromSubcollection(ctx context.Context, sheetID string) ([]*domain.SheetMember, error) {
	iter := r.collection.Doc(sheetID).Collection("members").Documents(ctx)
	defer iter.Stop()

	members := make([]*domain.SheetMember, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("iterate members: %w", err)
		}

		var member domain.SheetMember
		if err := doc.DataTo(&member); err != nil {
			return nil, fmt.Errorf("unmarshal member: %w", err)
		}

		member.SheetID = sheetID
		member.UserID = doc.Ref.ID
		members = append(members, &member)
	}

	return members, nil
}
