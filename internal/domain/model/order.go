package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type (
	Order struct {
		Id                uuid.UUID      `json:"order_uid"`
		TrackNumber       string         `json:"track_number"`
		Entry             string         `json:"entry"`
		DateCreated       time.Time      `json:"date_created"`
		Locale            sql.NullString `json:"locale"`
		InternalSignature sql.NullString `json:"internal_signature"`
		CustomerId        sql.NullString `json:"customer_id"`
		DeliveryService   sql.NullString `json:"delivery_service"`
		Shardkey          sql.NullString `json:"shardkey"`
		SmId              sql.NullInt64  `json:"sm_id"`
		OofShard          sql.NullString `json:"oof_shard"`
	}
	CompleteOrder struct {
		*Order
		*Delivery `json:"delivery"`
		*Payment  `json:"payment"`
		Items    []*Item `json:"items"`
	}
)
