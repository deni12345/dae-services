package sheet

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-core/internal/domain"
	"google.golang.org/api/iterator"
)

const (
	MemberRoleHost   = "host"
	MemberRoleMember = "member"
)

// AddMember adds a user to a sheet's member list using both denormalized array and subcollection
func (r *sheetRepo) AddMember(ctx context.Context, sheetID, userID string) error {
	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		sheetRef := r.collection.Doc(sheetID)
		snap, err := tx.Get(sheetRef)
		if err != nil {
			return fmt.Errorf("get sheet: %w", err)
		}

		var sheet domain.Sheet
		if err := snap.DataTo(&sheet); err != nil {
			return fmt.Errorf("unmarshal sheet: %w", err)
		}

		// Check if already a member
		for _, id := range sheet.MemberIDs {
			if id == userID {
				return nil // already member, idempotent
			}
		}

		// Update denormalized array
		sheet.MemberIDs = append(sheet.MemberIDs, userID)
		now := time.Now().UTC()

		updates := []firestore.Update{
			{Path: "member_ids", Value: sheet.MemberIDs},
			{Path: "updated_at", Value: now},
		}

		if err := tx.Update(sheetRef, updates); err != nil {
			return fmt.Errorf("update sheet members: %w", err)
		}

		// Sync to subcollection for efficient querying
		memberRef := sheetRef.Collection("members").Doc(userID)
		memberData := map[string]interface{}{
			"user_id":   userID, // Add for CollectionGroup query optimization
			"role":      MemberRoleMember,
			"joined_at": now,
		}

		return tx.Set(memberRef, memberData)
	})
}

// RemoveMember removes a user from a sheet's member list
func (r *sheetRepo) RemoveMember(ctx context.Context, sheetID, userID string) error {
	return r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
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
}

// ListMemberIDs returns all member IDs for a sheet
func (r *sheetRepo) ListMemberIDs(ctx context.Context, sheetID string) ([]string, error) {
	snap, err := r.collection.Doc(sheetID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get sheet: %w", err)
	}

	var sheet domain.Sheet
	if err := snap.DataTo(&sheet); err != nil {
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
