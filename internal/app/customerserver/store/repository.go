package store

import (
	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/model"
	"github.com/google/uuid"
)

type OfficeRepository interface {
	Add(*model.Office) error
	GetList() ([]*model.Office, error)
	Find(uuid.UUID) (*model.Office, error)
}

type UserRepository interface {
	Add(*model.User) error
	GetList() ([]*model.User, error)
}
