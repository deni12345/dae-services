package domain

import "time"

type Sheet struct {
	ID          string   `firestore:"-" json:"id"`
	Name        string   `firestore:"name" json:"name"`
	Description string   `firestore:"description" json:"description"`
	HostUserID  string   `firestore:"host_user_id"  json:"host_user_id"`
	Status      Status   `firestore:"status"         json:"status"`
	DeliveryFee Money    `firestore:"delivery_fee"   json:"delivery_fee"`
	Discount    int32    `firestore:"discount"       json:"discount"`
	MemberIDs   []string `firestore:"member_ids" json:"member_ids"` // Denormalized for backward compat

	// Optimistic locking / auditing
	UpdatedAt time.Time `firestore:"updated_at" json:"updated_at"`
	CreatedAt time.Time `firestore:"created_at" json:"created_at"`
}

func (s *Sheet) IsOpen() bool { return s.Status == Status_OPEN }

// SheetMember represents membership in sheets/{sheetID}/members/{userID} subcollection
type SheetMember struct {
	SheetID  string    `firestore:"-"`
	UserID   string    `firestore:"-"`
	Role     string    `firestore:"role" json:"role"` // "host" | "member"
	JoinedAt time.Time `firestore:"joined_at" json:"joined_at"`
}
