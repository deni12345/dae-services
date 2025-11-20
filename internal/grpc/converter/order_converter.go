package converter

import (
	"github.com/deni12345/dae-core/internal/app/order"
	"github.com/deni12345/dae-core/internal/domain"
	corev1 "github.com/deni12345/dae-core/proto/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Proto to DTO conversions

func CreateOrderReqFromProto(req *corev1.CreateOrderReq) *order.CreateOrderReq {
	lines := make([]order.OrderLineReq, len(req.Lines))
	for i, line := range req.Lines {
		options := make([]order.OrderLineOptionReq, len(line.Options))
		for j, opt := range line.Options {
			options[j] = order.OrderLineOptionReq{
				GroupID:  opt.GroupId,
				OptionID: opt.OptionId,
				Quantity: int(opt.Quantity),
			}
		}
		lines[i] = order.OrderLineReq{
			MenuItemID: line.MenuItemId,
			Options:    options,
			Quantity:   int(line.Quantity),
			Note:       line.Note,
		}
	}

	return &order.CreateOrderReq{
		SheetID: req.SheetId,
		Lines:   lines,
		Note:    req.Note,
		UserID:  req.UserId,
	}
}

func UpdateOrderReqFromProto(req *corev1.UpdateOrderReq) *order.UpdateOrderReq {
	lines := make([]order.OrderLineReq, len(req.Lines))
	for i, line := range req.Lines {
		options := make([]order.OrderLineOptionReq, len(line.Options))
		for j, opt := range line.Options {
			options[j] = order.OrderLineOptionReq{
				GroupID:  opt.GroupId,
				OptionID: opt.OptionId,
				Quantity: int(opt.Quantity),
			}
		}
		lines[i] = order.OrderLineReq{
			MenuItemID: line.MenuItemId,
			Options:    options,
			Quantity:   int(line.Quantity),
			Note:       line.Note,
		}
	}

	return &order.UpdateOrderReq{
		ID:    req.GetId(),
		Lines: lines,
		Note:  req.GetNote(),
	}
}

func ListOrdersReqFromProto(req *corev1.ListOrdersReq) *order.ListOrdersReq {
	dto := &order.ListOrdersReq{
		Limit: req.GetPageSize(),
	}

	// Parse cursor if provided
	if cursor := req.GetCursor(); cursor != nil && cursor.GetId() != "" {
		dto.Cursor = cursor.GetId()
	}

	// Parse filter if provided
	if filter := req.GetFilter(); filter != nil {
		dto.SheetID = filter.GetSheetId()

		if userID := filter.GetUserId(); userID != "" {
			dto.Filter.UserID = &userID
		}

		if since := filter.GetSince(); since != nil {
			t := since.AsTime()
			dto.Filter.Since = &t
		}
	}

	return dto
} // Domain to Proto conversions

func OrderToProto(o *domain.Order) *corev1.Order {
	if o == nil {
		return nil
	}

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

func OrderLineToProto(line domain.OrderLine) *corev1.OrderLine {
	options := make([]*corev1.OrderLineOption, len(line.Options))
	for i, opt := range line.Options {
		options[i] = &corev1.OrderLineOption{
			GroupId:    opt.GroupID,
			OptionId:   opt.OptionID,
			Title:      opt.Title,
			PriceDelta: MoneyToProto(opt.PriceDelta),
			Quantity:   opt.Quantity,
		}
	}

	return &corev1.OrderLine{
		MenuItemId:        line.MenuItemID,
		Name:              line.Name,
		Quantity:          line.Quantity,
		OrderBasePrice:    MoneyToProto(line.OrderBasePrice),
		OrderOptionsTotal: MoneyToProto(line.OrderOptionsTotal),
		OrderTotal:        MoneyToProto(line.OrderTotal),
		Options:           options,
		Note:              line.Note,
	}
}

func MoneyToProto(m domain.Money) *corev1.Money {
	return &corev1.Money{
		CurrencyCode: m.CurrencyCode,
		Amount:       m.Amount,
	}
}

func OrdersToProto(orders []*domain.Order) []*corev1.Order {
	result := make([]*corev1.Order, len(orders))
	for i, o := range orders {
		result[i] = OrderToProto(o)
	}
	return result
}

// ListOrdersRespToProto converts DTO response to proto
func ListOrdersRespToProto(resp *order.ListOrdersResp) *corev1.ListOrdersResp {
	if resp == nil {
		return &corev1.ListOrdersResp{}
	}

	protoResp := &corev1.ListOrdersResp{
		Orders: OrdersToProto(resp.Orders),
	}

	// Add next cursor if available
	if resp.NextCursor != "" {
		protoResp.NextCursor = &corev1.Cursor{
			Id: resp.NextCursor,
		}
	}

	return protoResp
}
