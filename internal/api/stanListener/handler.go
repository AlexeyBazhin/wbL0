package stanListener

import (
	"encoding/json"

	"github.com/AlexeyBazhin/wbL0/internal/api"
	"github.com/AlexeyBazhin/wbL0/internal/domain/model"
	"github.com/nats-io/stan.go"
)

type (
	Data struct {
		Type          string        `json:"type"`
		RecievedModel api.ModelJSON `json:"model"`
	}
	// CompleteModel struct {
	// 	*model.Order
	// 	*model.Delivery
	// 	*model.Payment
	// 	Items []*model.Item
	// }
)

func (stanListener *StanListener) stanHandler(msg *stan.Msg) {
	stanListener.logger.Info("[stan-listener] DATA RECEIVED")
	recievedData := &Data{}
	if err := json.Unmarshal(msg.Data, recievedData); err != nil {
		stanListener.logger.Errorf("[stan-listener] cannot unmarshal msg: %w", err)
		return
	}
	if recievedData.Type != "order" {
		stanListener.logger.Error("[stan-listener] invalid data type: %w")
		return
	}
	recievedOrder := api.OrderJSON{
		OrderUid:          recievedData.RecievedModel.OrderUid,
		TrackNumber:       recievedData.RecievedModel.TrackNumber,
		Entry:             recievedData.RecievedModel.Entry,
		Locale:            recievedData.RecievedModel.Locale,
		InternalSignature: recievedData.RecievedModel.InternalSignature,
		CustomerId:        recievedData.RecievedModel.CustomerId,
		DeliveryService:   recievedData.RecievedModel.DeliveryService,
		Shardkey:          recievedData.RecievedModel.Shardkey,
		SmId:              recievedData.RecievedModel.SmId,
		DateCreated:       recievedData.RecievedModel.DateCreated,
		OofShard:          recievedData.RecievedModel.OofShard,
	}
	order, err := stanListener.svc.CreateOrder(stanListener.ctx, recievedOrder)
	if err != nil {
		stanListener.logger.Errorf("[stan-listener] failed to create order: %w", err)
		return
	}
	delivery, err := stanListener.svc.CreateDelivery(stanListener.ctx, recievedData.RecievedModel.DeliveryJSON, recievedData.RecievedModel.OrderUid)
	if err != nil {
		stanListener.logger.Errorf("[stan-listener] failed to create delivery: %w", err)
		return
	}
	payment, err := stanListener.svc.CreatePayment(stanListener.ctx, recievedData.RecievedModel.PaymentJSON, recievedData.RecievedModel.OrderUid)
	if err != nil {
		stanListener.logger.Errorf("[stan-listener] failed to create payment: %w", err)
		return
	}
	items := make([]*model.Item, len(recievedData.RecievedModel.Items))
	for i, recievedItem := range recievedData.RecievedModel.Items {
		if items[i], err = stanListener.svc.CreateItem(stanListener.ctx, recievedItem, recievedData.RecievedModel.OrderUid); err != nil {
			stanListener.logger.Errorf("[stan-listener] failed to create item %v: %w", i, err)
			return
		}
	}

	completeModel := &model.CompleteOrder{
		Order:    order,
		Delivery: delivery,
		Payment:  payment,
		Items:    items,
	}
	if err := stanListener.svc.InsertCompleteOrder(stanListener.ctx, completeModel); err != nil {
		stanListener.logger.Errorf("[stan-listener] failed to insert complete order: %w", err)
		return
	}
	stanListener.logger.Info("[stan-listener] DATA SAVED TO PG")

	stanListener.PushToRedis(completeModel)
}

func (stanListener *StanListener) PushToRedis(completeModel *model.CompleteOrder) {
	modelJSON := api.MakeJSONModel(completeModel)
	modelByte, err := json.Marshal(modelJSON)
	if err != nil {
		stanListener.logger.Errorf("[stan-listener] failed to marshal modelJSON")
		return
	}
	if err := stanListener.redisClient.
		Set(stanListener.ctx, completeModel.Order.Id.String(), modelByte, 0).
		Err(); err != nil {
		stanListener.logger.Errorf("[stan-listener] failed to save in Redis")
		return
	}
	stanListener.logger.Info("[stan-listener] DATA SAVED TO REDIS")
}
