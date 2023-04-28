package service

import (
	"context"

	"github.com/google/uuid"
)

func (svc *service) PushToCache(ctx context.Context, orderUid uuid.UUID, data []byte) error {
	return svc.cache.PushToCache(ctx, orderUid.String(), data)
}

func (svc *service) PullFromCache(ctx context.Context, orderUid uuid.UUID) ([]byte, error) {
	return svc.cache.PullFromCache(ctx, orderUid.String())
}
