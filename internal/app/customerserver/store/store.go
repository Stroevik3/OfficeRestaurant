package store

type Store interface {
	Office() OfficeRepository
	User() UserRepository
}
