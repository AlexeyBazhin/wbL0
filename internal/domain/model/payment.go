package model

import (
	"github.com/google/uuid"
)

type (
	Payment struct {
		Id           int
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
)
