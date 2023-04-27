package api

import (
	"time"

	"github.com/AlexeyBazhin/wbL0/internal/domain/model"
	"github.com/google/uuid"
)

type (
	ErrorJSON struct {
		Err string `json:"error"`
	}
	ModelJSON struct {
		OrderUid          uuid.UUID `json:"order_uid"`
		TrackNumber       string    `json:"track_number"`
		Entry             string    `json:"entry"`
		DeliveryJSON      `json:"delivery"`
		PaymentJSON       `json:"payment"`
		Items             []ItemJSON `json:"items"`
		Locale            string     `json:"locale"`
		InternalSignature string     `json:"internal_signature"`
		CustomerId        string     `json:"customer_id"`
		DeliveryService   string     `json:"delivery_service"`
		Shardkey          string     `json:"shardkey"`
		SmId              int        `json:"sm_id"`
		DateCreated       time.Time  `json:"date_created"`
		OofShard          string     `json:"oof_shard"`
	}
	DeliveryJSON struct {
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		Zip     string `json:"zip"`
		City    string `json:"city"`
		Address string `json:"address"`
		Region  string `json:"region"`
		Email   string `json:"email"`
	}
	PaymentJSON struct {
		Transaction  uuid.UUID `json:"transaction"`
		RequestId    string    `json:"request_id"`
		Currency     string    `json:"currency"`
		Provider     string    `json:"provider"`
		Amount       int       `json:"amount"`
		PaymentDt    int       `json:"payment_dt"`
		Bank         string    `json:"bank"`
		DeliveryCost int       `json:"delivery_cost"`
		GoodsTotal   int       `json:"goods_total"`
		CustomFee    int       `json:"custom_fee"`
	}
	ItemJSON struct {
		ChrtId      int    `json:"chrt_id"`
		TrackNumber string `json:"track_number"`
		Price       int    `json:"price"`
		RId         string `json:"rid"`
		Name        string `json:"name"`
		Sale        int    `json:"sale"`
		Size        string `json:"size"`
		TotalPrice  int    `json:"total_price"`
		NmId        int    `json:"nm_id"`
		Brand       string `json:"brand"`
		Status      int    `json:"status"`
	}
	OrderJSON struct {
		OrderUid          uuid.UUID `json:"order_uid"`
		TrackNumber       string    `json:"track_number"`
		Entry             string    `json:"entry"`
		Locale            string    `json:"locale"`
		InternalSignature string    `json:"internal_signature"`
		CustomerId        string    `json:"customer_id"`
		DeliveryService   string    `json:"delivery_service"`
		Shardkey          string    `json:"shardkey"`
		SmId              int       `json:"sm_id"`
		DateCreated       time.Time `json:"date_created"`
		OofShard          string    `json:"oof_shard"`
	}
)

func MakeJSONModel(completeOrder *model.CompleteOrder) ModelJSON {
	items := make([]ItemJSON, len(completeOrder.Items))
	for i, item := range completeOrder.Items {
		items[i] = ItemJSON{
			ChrtId:      item.ChrtId,
			TrackNumber: item.TrackNumber,
			Price:       item.Price,
			RId:         item.RId,
			Name:        item.Name,
			Sale:        item.Sale,
			Size:        item.Size,
			TotalPrice:  item.TotalPrice,
			NmId:        item.NmId,
			Brand:       item.Brand,
			Status:      item.Status,
		}
	}
	modelJSON := ModelJSON{
		OrderUid:          completeOrder.Order.Id,
		TrackNumber:       completeOrder.Order.TrackNumber,
		Entry:             completeOrder.Order.Entry,
		Locale:            completeOrder.Order.Locale.String,
		InternalSignature: completeOrder.Order.InternalSignature.String,
		CustomerId:        completeOrder.Order.CustomerId.String,
		DeliveryService:   completeOrder.Order.DeliveryService.String,
		Shardkey:          completeOrder.Order.Shardkey.String,
		SmId:              int(completeOrder.Order.SmId.Int64),
		DateCreated:       completeOrder.Order.DateCreated,
		OofShard:          completeOrder.Order.OofShard.String,
		DeliveryJSON: DeliveryJSON{
			Name:    completeOrder.Delivery.Name,
			Phone:   completeOrder.Delivery.Phone,
			Zip:     completeOrder.Delivery.Zip,
			City:    completeOrder.Delivery.City,
			Address: completeOrder.Delivery.Address,
			Region:  completeOrder.Delivery.Region,
			Email:   completeOrder.Delivery.Email,
		},
		PaymentJSON: PaymentJSON{
			Transaction:  completeOrder.Payment.Transaction,
			RequestId:    completeOrder.Payment.RequestId,
			Currency:     completeOrder.Payment.Currency,
			Provider:     completeOrder.Payment.Provider,
			Amount:       completeOrder.Payment.Amount,
			PaymentDt:    completeOrder.Payment.PaymentDt,
			Bank:         completeOrder.Payment.Bank,
			DeliveryCost: completeOrder.Payment.DeliveryCost,
			GoodsTotal:   completeOrder.Payment.GoodsTotal,
			CustomFee:    completeOrder.Payment.CustomFee,
		},
		Items: items,
	}
	return modelJSON
}
