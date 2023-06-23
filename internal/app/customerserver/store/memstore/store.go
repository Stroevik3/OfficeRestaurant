package memstore

import (
	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/model"
	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/store"
	"github.com/google/uuid"
)

type Store struct {
	officeRepository *OfficeRepository
	userRepository   *UserRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) Office() store.OfficeRepository {
	if s.officeRepository != nil {
		return s.officeRepository
	}

	s.officeRepository = &OfficeRepository{
		store:   s,
		offices: make(map[uuid.UUID]*model.Office),
	}

	return s.officeRepository
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[uuid.UUID]*model.User),
	}

	return s.userRepository
}
