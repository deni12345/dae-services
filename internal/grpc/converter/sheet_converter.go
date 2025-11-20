package converter

import (
	"github.com/deni12345/dae-core/internal/app/sheet"
	"github.com/deni12345/dae-core/internal/domain"
	corev1 "github.com/deni12345/dae-core/proto/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Status mappings
var protoToDomainStatusMap = map[corev1.SheetStatus]domain.Status{
	corev1.SheetStatus_SHEET_STATUS_OPEN:        domain.Status_OPEN,
	corev1.SheetStatus_SHEET_STATUS_CLOSED:      domain.Status_CLOSED,
	corev1.SheetStatus_SHEET_STATUS_PENDING:     domain.Status_PENDING,
	corev1.SheetStatus_SHEET_STATUS_UNSPECIFIED: domain.Status_UNKNOWN,
}

var domainToProtoStatusMap = map[domain.Status]corev1.SheetStatus{
	domain.Status_OPEN:    corev1.SheetStatus_SHEET_STATUS_OPEN,
	domain.Status_CLOSED:  corev1.SheetStatus_SHEET_STATUS_CLOSED,
	domain.Status_PENDING: corev1.SheetStatus_SHEET_STATUS_PENDING,
	domain.Status_UNKNOWN: corev1.SheetStatus_SHEET_STATUS_UNSPECIFIED,
}

// CreateSheetReqFromProto converts proto CreateSheetReq to DTO
func CreateSheetReqFromProto(req *corev1.CreateSheetReq) *sheet.CreateSheetReq {
	var deliveryFee *domain.Money
	if protoFee := req.GetDeliveryFee(); protoFee != nil {
		deliveryFee = &domain.Money{
			CurrencyCode: protoFee.GetCurrencyCode(),
			Amount:       protoFee.GetAmount(),
		}
	}

	return &sheet.CreateSheetReq{
		IdempotencyKey: req.GetIdempotencyKey(),
		Name:           req.GetName(),
		Description:    req.GetDescription(),
		HostUserID:     req.GetHostUserId(),
		DeliveryFee:    deliveryFee,
		Discount:       req.GetDiscount(),
		MemberIDs:      req.GetMemberIds(),
		MenuItems:      MenuItemsFromProto(req.GetItems()),
	}
}

// MenuItemsFromProto converts proto MenuItems to DTO MenuItemReq
func MenuItemsFromProto(protoItems []*corev1.MenuItem) []sheet.MenuItemReq {
	if len(protoItems) == 0 {
		return nil
	}

	items := make([]sheet.MenuItemReq, len(protoItems))
	for i, item := range protoItems {
		var price int64
		var currency string
		if p := item.GetPrice(); p != nil {
			price = p.GetAmount()
			currency = p.GetCurrencyCode()
		}

		items[i] = sheet.MenuItemReq{
			ID:           item.GetId(),
			Name:         item.GetTitle(),
			Description:  item.GetDescription(),
			Active:       item.GetAvailable(),
			Price:        price,
			Currency:     currency,
			OptionGroups: MenuOptionGroupsFromProto(item.GetOptionGroups()),
		}
	}
	return items
}

// MenuOptionGroupsFromProto converts proto MenuOptionGroups to DTO
func MenuOptionGroupsFromProto(protoGroups []*corev1.MenuOptionGroup) []sheet.MenuOptionGroupReq {
	if len(protoGroups) == 0 {
		return nil
	}

	groups := make([]sheet.MenuOptionGroupReq, len(protoGroups))
	for i, grp := range protoGroups {
		groups[i] = sheet.MenuOptionGroupReq{
			ID:          grp.GetId(),
			Name:        grp.GetTitle(),
			Required:    grp.GetRequired(),
			MultiSelect: grp.GetMultiSelect(),
			MinSelect:   grp.GetMinSelect(),
			MaxSelect:   grp.GetMaxSelect(),
			Options:     MenuOptionsFromProto(grp.GetOptions()),
		}
	}
	return groups
}

// MenuOptionsFromProto converts proto MenuOptions to DTO
func MenuOptionsFromProto(protoOptions []*corev1.MenuOption) []sheet.MenuOptionReq {
	if len(protoOptions) == 0 {
		return nil
	}

	options := make([]sheet.MenuOptionReq, len(protoOptions))
	for i, opt := range protoOptions {
		var price int64
		if p := opt.GetPriceDelta(); p != nil {
			price = p.GetAmount()
		}

		options[i] = sheet.MenuOptionReq{
			ID:     opt.GetId(),
			Name:   opt.GetTitle(),
			Price:  price,
			Active: opt.GetAvailable(),
		}
	}
	return options
}

// SheetToProto converts domain Sheet to proto
func SheetToProto(s *domain.Sheet) *corev1.Sheet {
	if s == nil {
		return nil
	}

	return &corev1.Sheet{
		Id:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		HostUserId:  s.HostUserID,
		DeliveryFee: &corev1.Money{
			CurrencyCode: s.DeliveryFee.CurrencyCode,
			Amount:       s.DeliveryFee.Amount,
		},
		Discount:  s.Discount,
		Status:    domainToProtoStatusMap[s.Status],
		CreatedAt: timestamppb.New(s.CreatedAt),
		UpdatedAt: timestamppb.New(s.UpdatedAt),
	}
}

// UpdateSheetReqFromProto converts proto UpdateSheetReq to DTO
func UpdateSheetReqFromProto(req *corev1.UpdateSheetReq) *sheet.UpdateSheetReq {
	dto := &sheet.UpdateSheetReq{
		ID: req.GetId(),
	}

	if req.Name != nil {
		dto.Name = req.Name
	}

	if req.Description != nil {
		dto.Description = req.Description
	}

	if req.Status != nil {
		status := protoToDomainStatusMap[*req.Status]
		dto.Status = &status
	}

	return dto
}

// JoinSheetReqFromProto converts proto JoinSheetRequest to DTO
func JoinSheetReqFromProto(req *corev1.JoinSheetRequest) *sheet.JoinSheetReq {
	return &sheet.JoinSheetReq{
		SheetID: req.GetSheetId(),
		UserID:  req.GetUserId(),
	}
}

// SheetMemberToProto converts domain SheetMember to proto
func SheetMemberToProto(m *domain.SheetMember) *corev1.SheetMember {
	if m == nil {
		return nil
	}

	return &corev1.SheetMember{
		UserId:   m.UserID,
		SheetId:  m.SheetID,
		JoinedAt: timestamppb.New(m.JoinedAt),
	}
}

// ListSheetsReqFromProto converts proto ListSheetsReq to DTO
func ListSheetsReqFromProto(req *corev1.ListSheetsReq) *sheet.ListSheetsReq {
	dto := &sheet.ListSheetsReq{
		Limit: int(req.GetPageSize()),
	}

	// Extract cursor if provided (ID only)
	if cursor := req.GetCursor(); cursor != nil && cursor.GetId() != "" {
		dto.Cursor = cursor.GetId()
	}

	// Extract filter if provided
	if filter := req.GetFilter(); filter != nil {
		if ownerID := filter.GetOwnerUserId(); ownerID != "" {
			dto.HostUserID = &ownerID
		}
	}

	return dto
}

// ListSheetsRespToProto converts DTO ListSheetsResp to proto
func ListSheetsRespToProto(resp *sheet.ListSheetsResp) *corev1.ListSheetsResp {
	if resp == nil {
		return &corev1.ListSheetsResp{}
	}

	sheets := make([]*corev1.Sheet, len(resp.Sheets))
	for i, s := range resp.Sheets {
		sheets[i] = SheetToProto(s)
	}

	protoResp := &corev1.ListSheetsResp{
		Sheets: sheets,
	}

	// Add next cursor if available (ID only)
	if resp.NextCursor != "" {
		protoResp.NextCursor = &corev1.Cursor{
			Id: resp.NextCursor,
		}
	}

	return protoResp
}

// SheetsToProto converts slice of domain Sheets to proto
func SheetsToProto(sheets []*domain.Sheet) []*corev1.Sheet {
	result := make([]*corev1.Sheet, len(sheets))
	for i, s := range sheets {
		result[i] = SheetToProto(s)
	}
	return result
}
