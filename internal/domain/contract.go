package domain

import (
	"context"

	"github.com/AlexeyBazhin/wbL0/internal/domain/model"
	"github.com/google/uuid"
)

type (
	Repository interface {
		GetOrderModelById(ctx context.Context, orderUid uuid.UUID) (*model.CompleteOrder, error)
		InsertOrder(ctx context.Context, order *model.Order) error
		InsertDelivery(ctx context.Context, order *model.Delivery) error
		InsertPayment(ctx context.Context, order *model.Payment) error
		InsertItem(ctx context.Context, order *model.Item) error
	}
	Cache interface {
		PushToCache(ctx context.Context, orderUidStr string, data []byte) error
		PullFromCache(ctx context.Context, orderUidStr string) ([]byte, error)
	}
)
