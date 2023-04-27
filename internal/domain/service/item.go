package service

import (
	"context"

	"github.com/AlexeyBazhin/wbL0/internal/api"
	"github.com/AlexeyBazhin/wbL0/internal/domain/model"
	"github.com/google/uuid"
)

func (svc *service) CreateItem(ctx context.Context, recievedItem api.ItemJSON, orderUid uuid.UUID) (*model.Item, error) {
	return &model.Item{
		OrderId:     orderUid,
		ChrtId:      recievedItem.ChrtId,
		TrackNumber: recievedItem.TrackNumber,
		Price:       recievedItem.Price,
		RId:         recievedItem.RId,
		Name:        recievedItem.Name,
		Sale:        recievedItem.Sale,
		Size:        recievedItem.Size,
		TotalPrice:  recievedItem.TotalPrice,
		NmId:        recievedItem.NmId,
		Brand:       recievedItem.Brand,
		Status:      recievedItem.Status,
	}, nil
}
