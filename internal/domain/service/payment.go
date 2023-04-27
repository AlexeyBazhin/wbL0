package service

import (
	"context"

	"github.com/AlexeyBazhin/wbL0/internal/api"
	"github.com/AlexeyBazhin/wbL0/internal/domain/model"
	"github.com/google/uuid"
)

func (svc *service) CreatePayment(ctx context.Context, recievedPayment api.PaymentJSON, orderUid uuid.UUID) (*model.Payment, error) {
	return &model.Payment{
		Transaction:  orderUid,
		RequestId:    recievedPayment.RequestId,
		Currency:     recievedPayment.Currency,
		Provider:     recievedPayment.Provider,
		Amount:       recievedPayment.Amount,
		PaymentDt:    recievedPayment.PaymentDt,
		Bank:         recievedPayment.Bank,
		DeliveryCost: recievedPayment.DeliveryCost,
		GoodsTotal:   recievedPayment.GoodsTotal,
		CustomFee:    recievedPayment.CustomFee,
	}, nil
}
