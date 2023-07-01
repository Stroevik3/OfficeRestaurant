package pgstore

import (
	"database/sql"
	"time"

	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/model"
	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/store"
	"github.com/google/uuid"
)

type OfficeRepository struct {
	store *Store
}

func (r *OfficeRepository) Add(o *model.Office) error {
	var err error
	o.Uuid, err = uuid.NewUUID()
	if err != nil {
		return err
	}
	r.store.db.QueryRow(
		"INSERT INTO offices (id, name, addres, createdAt) VALUES ($1, $2, $3, $4)",
		o.Uuid, o.Name, o.Addres, time.Now(),
	)
	return nil
}

func (r *OfficeRepository) GetList() ([]*model.Office, error) {
	rows, err := r.store.db.Query("SELECT id, name, addres, createdAt FROM offices")
	if err != nil {
		return nil, err
	}
	var newOffices []*model.Office
	for rows.Next() {
		e := new(model.Office)
		err := rows.Scan(&e.Uuid, &e.Name, &e.Addres, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		newOffices = append(newOffices, e)
	}

	return newOffices, nil
}

func (r *OfficeRepository) Find(id uuid.UUID) (*model.Office, error) {
	of := &model.Office{}
	if err := r.store.db.QueryRow(
		"SELECT id, name, addres, createdAt FROM offices WHERE id = $1",
		id,
	).Scan(
		&of.Uuid,
		&of.Name,
		&of.Addres,
		&of.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return of, nil
}
