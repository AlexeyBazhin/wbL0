package repository

import (
	"context"
	"fmt"

	"github.com/AlexeyBazhin/wbL0/internal/domain/model"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (repo *Repository) InsertOrder(ctx context.Context, order *model.Order) error {
	_, err := repo.sq.Insert("orders").
		Columns("order_uid", "track_number",
			"entry", "date_created",
			"locale", "internal_signature",
			"customer_id", "delivery_service",
			"shardkey", "sm_id",
			"oof_shard").
		Values(order.Id, order.TrackNumber,
			order.Entry, order.DateCreated,
			order.Locale, order.InternalSignature,
			order.CustomerId, order.DeliveryService,
			order.Shardkey, order.SmId,
			order.OofShard).ExecContext(ctx)
	return err
}

func (repo *Repository) GetOrderModelById(ctx context.Context, orderUid uuid.UUID) (*model.CompleteOrder, error) {
	row := repo.sq.Select("o.order_uid", "o.track_number",
		"o.entry", "o.date_created", "o.locale",
		"o.internal_signature", "o.customer_id",
		"o.delivery_service", "o.shardkey",
		"o.sm_id", "o.oof_shard",
		"d.delivery_id", "d.order_uid",
		"d.name", "d.phone", "d.zip",
		"d.city", "d.address", "d.region",
		"d.email",
		"p.payment_id", "p.transaction", "p.request_id",
		"p.currency", "p.provider", "p.amount",
		"p.payment_dt", "p.bank", "p.delivery_cost",
		"p.goods_total", "p.custom_fee").
		From("orders o").
		Join("deliveries d ON d.order_uid = o.order_uid").
		Join("payments p ON p.transaction = o.order_uid").
		Where(squirrel.Eq{"o.order_uid": orderUid}).
		QueryRowContext(ctx)

	orderModel := &model.CompleteOrder{
		Order:    &model.Order{},
		Delivery: &model.Delivery{},
		Payment:  &model.Payment{},
	}

	if err := row.Scan(
		&orderModel.Order.Id, &orderModel.Order.TrackNumber,
		&orderModel.Order.Entry, &orderModel.Order.DateCreated,
		&orderModel.Order.Locale, &orderModel.Order.InternalSignature,
		&orderModel.Order.CustomerId, &orderModel.Order.DeliveryService,
		&orderModel.Order.Shardkey, &orderModel.Order.SmId,
		&orderModel.Order.OofShard,
		&orderModel.Delivery.Id, &orderModel.Delivery.OrderId,
		&orderModel.Delivery.Name, &orderModel.Delivery.Phone,
		&orderModel.Delivery.Zip, &orderModel.Delivery.City,
		&orderModel.Delivery.Address, &orderModel.Delivery.Region,
		&orderModel.Delivery.Email,
		&orderModel.Payment.Id, &orderModel.Payment.Transaction,
		&orderModel.Payment.RequestId, &orderModel.Payment.Currency,
		&orderModel.Payment.Provider, &orderModel.Payment.Amount,
		&orderModel.Payment.PaymentDt, &orderModel.Payment.Bank,
		&orderModel.Payment.DeliveryCost, &orderModel.Payment.GoodsTotal,
		&orderModel.Payment.CustomFee,
	); err != nil {
		return nil, fmt.Errorf("error while performing sql request: %w", err)
	}

	rows, err := repo.sq.Select("i.item_id", "i.order_uid",
		"i.chrt_id", "i.track_number", "i.price",
		"i.rid", "i.name", "i.sale", "i.size",
		"i.total_price", "i.nm_id", "i.brand",
		"i.status").
		From("items i").
		Where(squirrel.Eq{"i.order_uid": orderUid}).
		QueryContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("error while performing sql request: %w", err)
	}

	defer func() {
		if err = rows.Close(); err != nil {
			repo.logger.Error("error while closing sql rows", zap.Error(err))
		}
	}()

	items := make([]*model.Item, 0)
	for rows.Next() {
		item := &model.Item{}
		if err = rows.Scan(
			&item.Id, &item.OrderId,
			&item.ChrtId, &item.TrackNumber,
			&item.Price, &item.RId, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice,
			&item.NmId, &item.Brand, &item.Status,
		); err != nil {
			return nil, fmt.Errorf("error while scanning sql row: %w", err)
		}
		items = append(items, item)
	}
	orderModel.Items = items
	return orderModel, nil
}
