package service

import (
	"context"

	"github.com/AlexeyBazhin/wbL0/internal/api"
	"github.com/AlexeyBazhin/wbL0/internal/domain/model"
	"github.com/google/uuid"
)

func (svc *service) GetOrderById(ctx context.Context, orderUid uuid.UUID) (*model.CompleteOrder, error) {
	return svc.repo.GetOrderModelById(ctx, orderUid)
}
func (svc *service) InsertCompleteOrder(ctx context.Context, completeModels *model.CompleteOrder) error {
	if err := svc.repo.InsertOrder(ctx, completeModels.Order); err != nil {
		return err
	}
	if err := svc.repo.InsertDelivery(ctx, completeModels.Delivery); err != nil {
		return err
	}
	if err := svc.repo.InsertPayment(ctx, completeModels.Payment); err != nil {
		return err
	}
	for _, item := range completeModels.Items {
		if err := svc.repo.InsertItem(ctx, item); err != nil {
			return err
		}
	}
	return nil
}
func (svc *service) CreateOrder(ctx context.Context, recievedOrder api.OrderJSON) (*model.Order, error) {
	order := &model.Order{
		Id:          recievedOrder.OrderUid,
		TrackNumber: recievedOrder.TrackNumber,
		Entry:       recievedOrder.Entry,
		DateCreated: recievedOrder.DateCreated,
		// Locale: recievedOrder.Locale,
		// InternalSignature: recievedOrder.InternalSignature,
		// CustomerId: recievedOrder.CustomerId,
		// DeliveryService: recievedOrder.DeliveryService,
		// Shardkey: recievedOrder.Shardkey,
		// SmId: recievedOrder.SmId,
		// OofShard: recievedOrder.OofShard,
	}
	order.Locale.Scan(recievedOrder.Locale)
	order.InternalSignature.Scan(recievedOrder.InternalSignature)
	order.CustomerId.Scan(recievedOrder.CustomerId)
	order.DeliveryService.Scan(recievedOrder.DeliveryService)
	order.Shardkey.Scan(recievedOrder.Shardkey)
	order.SmId.Scan(recievedOrder.SmId)
	order.OofShard.Scan(recievedOrder.OofShard)

	return order, nil
}

//ЕСЛИ БЫ БЫЛО РЕАЛИЗОВАНО СОЗДАНИЕ ЗАКАЗА ЧЕРЕЗ ПОСТ ЗАПРОС
// func (svc *service) CreateOrderFromRequest(orderReq server.OrderReq) (*model.Order, error) {
// 	//1 ВАРИАНТ
// 	//locale := &sql.NullString{}
// 	// if err := locale.Scan(req.Locale); err != nil {
// 	// 	return nil, err
// 	// }
// 	// internalSignature := &sql.NullString{}
// 	// if err := internalSignature.Scan(req.InternalSignature); err != nil {
// 	// 	return nil, err
// 	// }
// 	// customerId := &sql.NullString{}
// 	// if err := customerId.Scan(req.CustomerId); err != nil {
// 	// 	return nil, err
// 	// }
// 	// deliveryService := &sql.NullString{}
// 	// if err := deliveryService.Scan(req.DeliveryService); err != nil {
// 	// 	return nil, err
// 	// }
// 	// shardKey := &sql.NullString{}
// 	// if err := shardKey.Scan(req.Shardkey); err != nil {
// 	// 	return nil, err
// 	// }
// 	// smId := &sql.NullString{}
// 	// if err := smId.Scan(req.SmId); err != nil {
// 	// 	return nil, err
// 	// }
// 	// oofShard := &sql.NullString{}
// 	// if err := oofShard.Scan(req.OofShard); err != nil {
// 	// 	return nil, err
// 	// }
// 	// orderUid, err := uuid.NewRandom()
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// order := &model.Order{
// 	// 	Id:          orderUid,
// 	// 	TrackNumber: req.TrackNumber,
// 	// 	Entry:       req.Entry,
// 	// 	DateCreated: time.Now(),
// 	// 	Locale:      *locale,
// 	// 	InternalSignature: *internalSignature,
// 	// 	CustomerId: *customerId,
// 	// 	DeliveryService: *deliveryService,
// 	// 	Shardkey: *shardKey,
// 	// 	SmId: *smId,
// 	// 	OofShard: *oofShard,
// 	// }
// 	//2 ВАРИАНТ
// 	orderUid, err := uuid.NewRandom()
// 	if err != nil {
// 		return nil, err
// 	}
// 	order := &model.Order{
// 		Id:          orderUid,
// 		TrackNumber: orderReq.TrackNumber,
// 		Entry:       orderReq.Entry,
// 		DateCreated: time.Now(),
// 	}
// 	order.Locale.Scan(orderReq.Locale)
// 	order.InternalSignature.Scan(orderReq.InternalSignature)
// 	order.CustomerId.Scan(orderReq.CustomerId)
// 	order.DeliveryService.Scan(orderReq.DeliveryService)
// 	order.Shardkey.Scan(orderReq.Shardkey)
// 	order.SmId.Scan(orderReq.SmId)
// 	order.OofShard.Scan(orderReq.OofShard)
//		return order, svc.repo.CreateOrder(order)
//	}
