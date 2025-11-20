package domain

import (
	"time"

	corev1 "github.com/deni12345/dae-core/proto/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderStatus string

const (
	OrderStatusUnspecified OrderStatus = ""
	OrderStatusPending     OrderStatus = "pending"
	OrderStatusConfirmed   OrderStatus = "confirmed"
	OrderStatusCancelled   OrderStatus = "cancelled"
	OrderStatusCompleted   OrderStatus = "completed"
)

type OrderLineOption struct {
	GroupID    string `firestore:"group_id" json:"group_id"`
	OptionID   string `firestore:"option_id" json:"option_id"`
	Title      string `firestore:"title" json:"title"`
	PriceDelta Money  `firestore:"price_delta" json:"price_delta"`
	Quantity   int32  `firestore:"quantity" json:"quantity"`
}

type OrderLine struct {
	MenuItemID        string            `firestore:"menu_item_id" json:"menu_item_id"`
	Name              string            `firestore:"name" json:"name"`
	Quantity          int32             `firestore:"quantity" json:"quantity"`
	OrderBasePrice    Money             `firestore:"order_base_price" json:"order_base_price"`
	OrderOptionsTotal Money             `firestore:"order_options_total" json:"order_options_total"`
	OrderTotal        Money             `firestore:"order_total" json:"order_total"`
	Options           []OrderLineOption `firestore:"options" json:"options"`
	Note              string            `firestore:"note" json:"note"`
}

type Order struct {
	ID        string      `firestore:"-" json:"id"`
	SheetID   string      `firestore:"sheet_id" json:"sheet_id"`
	UserID    string      `firestore:"user_id" json:"user_id"`
	Lines     []OrderLine `firestore:"lines" json:"lines"`
	Subtotal  Money       `firestore:"subtotal" json:"subtotal"`
	Total     Money       `firestore:"total" json:"total"`
	Note      string      `firestore:"note" json:"note"`
	CreatedAt time.Time   `firestore:"created_at" json:"created_at"`
	UpdatedAt time.Time   `firestore:"updated_at" json:"updated_at"`
}

// GetMoneyAmount calculates the total amount in the smallest unit (considering nanos)
func (m Money) GetAmount() int64 {
	return m.Amount
}

// NewMoney creates a Money instance from amount in smallest units and currency
func NewMoney(amount int64, currencyCode string) Money {
	return Money{
		CurrencyCode: currencyCode,
		Amount:       amount,
	}
}

func CalculateOrderLineTotals(line *OrderLine) {
	optionsTotal := int64(0)
	for _, opt := range line.Options {
		optionsTotal += int64(opt.Quantity) * opt.PriceDelta.GetAmount()
	}

	line.OrderOptionsTotal = NewMoney(optionsTotal, line.OrderBasePrice.CurrencyCode)

	lineTotal := line.OrderBasePrice.GetAmount()*int64(line.Quantity) + optionsTotal
	line.OrderTotal = NewMoney(lineTotal, line.OrderBasePrice.CurrencyCode)
}

func CalculateOrderTotals(order *Order) {
	subtotal := int64(0)
	currency := ""
	for i := range order.Lines {
		CalculateOrderLineTotals(&order.Lines[i])
		subtotal += order.Lines[i].OrderTotal.GetAmount()
		if currency == "" {
			currency = order.Lines[i].OrderBasePrice.CurrencyCode
		}
	}

	order.Subtotal = NewMoney(subtotal, currency)
	order.Total = NewMoney(subtotal, currency)
}

// Conversion functions between domain and proto

func MoneyToProto(m Money) *corev1.Money {
	return &corev1.Money{
		CurrencyCode: m.CurrencyCode,
		Amount:       m.Amount,
	}
}

func MoneyFromProto(m *corev1.Money) Money {
	if m == nil {
		return Money{}
	}
	return Money{
		CurrencyCode: m.CurrencyCode,
		Amount:       m.Amount,
	}
}

func OrderLineOptionToProto(o OrderLineOption) *corev1.OrderLineOption {
	return &corev1.OrderLineOption{
		GroupId:    o.GroupID,
		OptionId:   o.OptionID,
		Title:      o.Title,
		PriceDelta: MoneyToProto(o.PriceDelta),
		Quantity:   o.Quantity,
	}
}

func OrderLineOptionFromProto(o *corev1.OrderLineOption) OrderLineOption {
	if o == nil {
		return OrderLineOption{}
	}
	return OrderLineOption{
		GroupID:    o.GroupId,
		OptionID:   o.OptionId,
		Title:      o.Title,
		PriceDelta: MoneyFromProto(o.PriceDelta),
		Quantity:   o.Quantity,
	}
}

func OrderLineToProto(l OrderLine) *corev1.OrderLine {
	options := make([]*corev1.OrderLineOption, len(l.Options))
	for i, opt := range l.Options {
		options[i] = OrderLineOptionToProto(opt)
	}

	return &corev1.OrderLine{
		MenuItemId:        l.MenuItemID,
		Name:              l.Name,
		Quantity:          l.Quantity,
		OrderBasePrice:    MoneyToProto(l.OrderBasePrice),
		OrderOptionsTotal: MoneyToProto(l.OrderOptionsTotal),
		OrderTotal:        MoneyToProto(l.OrderTotal),
		Options:           options,
		Note:              l.Note,
	}
}

func OrderLineFromProto(l *corev1.OrderLine) OrderLine {
	if l == nil {
		return OrderLine{}
	}

	options := make([]OrderLineOption, len(l.Options))
	for i, opt := range l.Options {
		options[i] = OrderLineOptionFromProto(opt)
	}

	return OrderLine{
		MenuItemID:        l.MenuItemId,
		Name:              l.Name,
		Quantity:          l.Quantity,
		OrderBasePrice:    MoneyFromProto(l.OrderBasePrice),
		OrderOptionsTotal: MoneyFromProto(l.OrderOptionsTotal),
		OrderTotal:        MoneyFromProto(l.OrderTotal),
		Options:           options,
		Note:              l.Note,
	}
}

func OrderToProto(o Order) *corev1.Order {
	lines := make([]*corev1.OrderLine, len(o.Lines))
	for i, line := range o.Lines {
		lines[i] = OrderLineToProto(line)
	}

	return &corev1.Order{
		Id:        o.ID,
		SheetId:   o.SheetID,
		UserId:    o.UserID,
		Lines:     lines,
		Subtotal:  MoneyToProto(o.Subtotal),
		Total:     MoneyToProto(o.Total),
		Note:      o.Note,
		CreateAt:  timestamppb.New(o.CreatedAt),
		UpdatedAt: timestamppb.New(o.UpdatedAt),
	}
}

func OrderFromProto(o *corev1.Order) Order {
	if o == nil {
		return Order{}
	}

	lines := make([]OrderLine, len(o.Lines))
	for i, line := range o.Lines {
		lines[i] = OrderLineFromProto(line)
	}

	createdAt := time.Time{}
	if o.CreateAt != nil {
		createdAt = o.CreateAt.AsTime()
	}

	updatedAt := time.Time{}
	if o.UpdatedAt != nil {
		updatedAt = o.UpdatedAt.AsTime()
	}

	return Order{
		ID:        o.Id,
		SheetID:   o.SheetId,
		UserID:    o.UserId,
		Lines:     lines,
		Subtotal:  MoneyFromProto(o.Subtotal),
		Total:     MoneyFromProto(o.Total),
		Note:      o.Note,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
