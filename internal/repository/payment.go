package repository

import (
	"context"
	"fmt"

	"github.com/AlexeyBazhin/wbL0/internal/domain/model"
)

func (repo *Repository) InsertPayment(ctx context.Context, payment *model.Payment) error {
	row := repo.sq.Insert("payments").
		Columns("transaction", "request_id",
			"currency", "provider",
			"amount", "payment_dt",
			"bank", "delivery_cost",
			"goods_total", "custom_fee").
		Values(payment.Transaction, payment.RequestId,
			payment.Currency, payment.Provider,
			payment.Amount, payment.PaymentDt,
			payment.Bank, payment.DeliveryCost,
			payment.GoodsTotal, payment.CustomFee).
		Suffix("RETURNING \"payment_id\"").
		QueryRowContext(ctx)

	if err := row.Scan(&payment.Id); err != nil {
		return fmt.Errorf("error while scanning sql row: %w", err)
	}
	return nil
}
