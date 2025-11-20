package order

import (
	"context"
	"fmt"

	"github.com/deni12345/dae-core/internal/domain"
)

// buildOrderLine calculates price, options, and builds OrderLine with groupIDs
func (u *usecase) buildOrderLine(ctx context.Context, sheetID string, lineReq OrderLineReq) (domain.OrderLine, error) {
	// Fetch menu item
	item, err := u.sheetRepo.GetMenuItemByID(ctx, sheetID, lineReq.MenuItemID)
	if err != nil {
		return domain.OrderLine{}, fmt.Errorf("get menu item: %w", err)
	}

	// Calculate pricing
	unitBase := item.Price
	var unitOpts int64
	var orderOptions []domain.OrderLineOption

	// Process options with group information
	if len(lineReq.Options) > 0 {
		for _, optGroup := range item.OptionGroups {
			for _, optReq := range lineReq.Options {
				if optReq.Quantity <= 0 {
					continue
				}

				optItem, ok := optGroup.Options[optReq.OptionID]
				if !ok {
					continue
				}

				// Validate single vs multi select
				if optReq.Quantity > 1 && optGroup.Type == domain.GroupSingle {
					return domain.OrderLine{}, fmt.Errorf("cannot select multiple options for single-select group: %s", optGroup.ID)
				}

				optionPrice := optItem.Price * int64(optReq.Quantity)
				unitOpts += optionPrice

				orderOptions = append(orderOptions, domain.OrderLineOption{
					GroupID:    optGroup.ID, // Include group ID
					OptionID:   optReq.OptionID,
					Title:      optItem.Name,
					PriceDelta: domain.NewMoney(optionPrice, item.Currency),
					Quantity:   int32(optReq.Quantity),
				})
			}
		}
	}

	unitTotal := unitBase + unitOpts
	lineTotal := unitTotal * int64(lineReq.Quantity)

	return domain.OrderLine{
		MenuItemID:        lineReq.MenuItemID,
		Name:              item.Name,
		Quantity:          int32(lineReq.Quantity),
		OrderBasePrice:    domain.NewMoney(unitBase, item.Currency),
		OrderOptionsTotal: domain.NewMoney(unitOpts, item.Currency),
		OrderTotal:        domain.NewMoney(lineTotal, item.Currency),
		Options:           orderOptions,
		Note:              lineReq.Note,
	}, nil
}
