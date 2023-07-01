package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Uuid              uuid.UUID
	Name, Description string
	Type              int
	Weight            int32
	Price             float32
	CreatedAt         time.Time
}
