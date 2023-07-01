package model

import "github.com/google/uuid"

type OrderItem struct {
	Count       int32
	ProductUuid uuid.UUID
}
