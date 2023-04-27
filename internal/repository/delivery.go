package repository

import (
	"context"
	"fmt"

	"github.com/AlexeyBazhin/wbL0/internal/domain/model"
)

func (repo *Repository) InsertDelivery(ctx context.Context, delivery *model.Delivery) error {
	row := repo.sq.Insert("deliveries").
		Columns("order_uid", "name",
			"phone", "zip",
			"city", "address",
			"region", "email").
		Values(delivery.OrderId, delivery.Name,
			delivery.Phone, delivery.Zip,
			delivery.City, delivery.Address,
			delivery.Region, delivery.Email).
		Suffix("RETURNING \"delivery_id\"").
		QueryRowContext(ctx)

	if err := row.Scan(&delivery.Id); err != nil {
		return fmt.Errorf("error while scanning sql row: %w", err)
	}
	return nil
}
