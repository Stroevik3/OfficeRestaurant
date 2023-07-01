package pgstore

import (
	"database/sql"

	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/store"
	_ "github.com/lib/pq"
)

type Store struct {
	db               *sql.DB
	officeRepository *OfficeRepository
	userRepository   *UserRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Office() store.OfficeRepository {
	if s.officeRepository != nil {
		return s.officeRepository
	}

	s.officeRepository = &OfficeRepository{
		store: s,
	}

	return s.officeRepository
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
