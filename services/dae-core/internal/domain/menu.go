package domain

type Per string

const (
	PerUnit  Per = "unit"
	PerOrder Per = "order"
)

type ItemID = string

type Option struct {
	ID     string `firestore:"id" json:"id"`
	Name   string `firestore:"name" json:"name"`
	Price  int64  `firestore:"price" json:"price"` // minor units
	Per    Per    `firestore:"per" json:"per"`     // unit|order
	Active bool   `firestore:"active" json:"active"`
}

type OptionGroupType string

const (
	GroupSingle OptionGroupType = "single"
	GroupMulti  OptionGroupType = "multi"
)

type OptionGroup struct {
	ID        string            `firestore:"id" json:"id"`
	Name      string            `firestore:"name" json:"name"`
	Type      OptionGroupType   `firestore:"type" json:"type"`
	Required  bool              `firestore:"required" json:"required"`
	MaxSelect int               `firestore:"max_select" json:"max_select"` // 0 = unlimited
	Options   map[string]Option `firestore:"options" json:"options"`
}

type MenuItem struct {
	ID           string                 `firestore:"id" json:"id"`
	Name         string                 `firestore:"name" json:"name"`
	Active       bool                   `firestore:"active" json:"active"`
	Price        int64                  `firestore:"price" json:"price"`
	Currency     string                 `firestore:"currency" json:"currency"`
	OptionGroups map[string]OptionGroup `firestore:"option_groups" json:"option_groups"`
	UpdatedAt    int64                  `firestore:"updated_at" json:"updated_at"` // unix seconds
}
