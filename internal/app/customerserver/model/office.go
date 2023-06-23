package model

import (
	"time"

	"github.com/google/uuid"
)

type Office struct {
	Uuid         uuid.UUID
	Name, Addres string
	CreatedAt    time.Time
}
