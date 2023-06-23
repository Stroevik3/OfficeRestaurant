package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Uuid       uuid.UUID
	Name       string
	OfficeUser *Office
	CreatedAt  time.Time
}
