package domain

type Money struct {
	CurrencyCode string `firestore:"currency_code" json:"currency_code"`
	Amount       int64  `firestore:"amount" json:"amount"`
}

type Status int32

const (
	Status_UNKNOWN Status = iota
	Status_OPEN
	Status_PENDING
	Status_CLOSED
)

var Status_name = map[Status]string{
	Status_UNKNOWN: "UNKNOWN",
	Status_OPEN:    "OPEN",
	Status_PENDING: "PAUSED",
	Status_CLOSED:  "CLOSED",
}
