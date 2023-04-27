package service

import (
	"context"

	"github.com/AlexeyBazhin/wbL0/internal/api"
	"github.com/AlexeyBazhin/wbL0/internal/domain/model"
	"github.com/google/uuid"
)

func (svc *service) CreateDelivery(ctx context.Context, recievedDelivery api.DeliveryJSON, orderUid uuid.UUID) (*model.Delivery, error) {
	return &model.Delivery{
		OrderId: orderUid,
		Name:    recievedDelivery.Name,
		Phone:   recievedDelivery.Phone,
		Zip:     recievedDelivery.Zip,
		City:    recievedDelivery.City,
		Address: recievedDelivery.Address,
		Region:  recievedDelivery.Region,
		Email:   recievedDelivery.Email,
	}, nil
}
