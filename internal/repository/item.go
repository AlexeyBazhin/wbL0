package repository

import (
	"context"
	"fmt"

	"github.com/AlexeyBazhin/wbL0/internal/domain/model"
)

func (repo *Repository) InsertItem(ctx context.Context, item *model.Item) error {
	row := repo.sq.Insert("items").
		Columns("order_uid", "chrt_id",
			"track_number", "price",
			"rid", "name",
			"sale", "size",
			"total_price", "nm_id",
			"brand", "status").
		Values(item.OrderId, item.ChrtId,
			item.TrackNumber, item.Price,
			item.RId, item.Name,
			item.Sale, item.Size,
			item.TotalPrice, item.NmId,
			item.Brand, item.Status).
		Suffix("RETURNING \"item_id\"").
		QueryRowContext(ctx)

	if err := row.Scan(&item.Id); err != nil {
		return fmt.Errorf("error while scanning sql row: %w", err)
	}
	return nil
}
