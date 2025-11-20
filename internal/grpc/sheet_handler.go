package grpc

import (
	"context"

	"github.com/deni12345/dae-core/internal/app/sheet"
	"github.com/deni12345/dae-core/internal/grpc/converter"
	"github.com/deni12345/dae-core/internal/grpc/errors"
	corev1 "github.com/deni12345/dae-core/proto/gen"
)

type SheetHandler struct {
	corev1.UnimplementedSheetsServiceServer
	uc sheet.Usecase
}

func NewSheetHandler(uc sheet.Usecase) *SheetHandler {
	return &SheetHandler{
		uc: uc,
	}
}

func (h *SheetHandler) CreateSheet(ctx context.Context, req *corev1.CreateSheetReq) (*corev1.CreateSheetResp, error) {
	sheet, err := h.uc.CreateSheet(ctx, converter.CreateSheetReqFromProto(req))
	if err != nil {
		return nil, errors.ToGRPCStatus(err)
	}

	return &corev1.CreateSheetResp{
		Sheet: converter.SheetToProto(sheet),
	}, nil
}

func (h *SheetHandler) GetSheet(ctx context.Context, req *corev1.GetSheetReq) (*corev1.GetSheetResp, error) {
	sheet, err := h.uc.GetSheet(ctx, req.Id)
	if err != nil {
		return nil, errors.ToGRPCStatus(err)
	}

	return &corev1.GetSheetResp{
		Sheet: converter.SheetToProto(sheet),
	}, nil
}

func (h *SheetHandler) UpdateSheet(ctx context.Context, req *corev1.UpdateSheetReq) (*corev1.UpdateSheetResp, error) {
	sheet, err := h.uc.UpdateSheet(ctx, converter.UpdateSheetReqFromProto(req))
	if err != nil {
		return nil, errors.ToGRPCStatus(err)
	}

	return &corev1.UpdateSheetResp{
		Sheet: converter.SheetToProto(sheet),
	}, nil
}

func (h *SheetHandler) ListSheets(ctx context.Context, req *corev1.ListSheetsReq) (*corev1.ListSheetsResp, error) {
	resp, err := h.uc.ListSheets(ctx, converter.ListSheetsReqFromProto(req))
	if err != nil {
		return nil, errors.ToGRPCStatus(err)
	}

	return converter.ListSheetsRespToProto(resp), nil
}
