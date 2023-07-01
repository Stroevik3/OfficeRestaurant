package pgstore

import (
	"time"

	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/model"
	"github.com/google/uuid"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Add(u *model.User) error {
	var err error
	u.Uuid, err = uuid.NewUUID()
	if err != nil {
		return err
	}
	r.store.db.QueryRow(
		"INSERT INTO users (id ,name, officeId, createdAt) VALUES ($1, $2, $3, $4)",
		u.Uuid, u.Name, u.OfficeUser.Uuid, time.Now(),
	)
	return nil
}

func (r *UserRepository) GetList() ([]*model.User, error) {
	rows, err := r.store.db.Query("SELECT u.id, u.name, u.createdAt, u.officeId, o.name, o.addres, o.createdAt FROM users u join offices o on o.id = u.officeId")
	if err != nil {
		return nil, err
	}
	var newUsers []*model.User
	for rows.Next() {
		u := new(model.User)
		u.OfficeUser = new(model.Office)
		err := rows.Scan(&u.Uuid, &u.Name, &u.CreatedAt, &u.OfficeUser.Uuid, &u.OfficeUser.Name, &u.OfficeUser.Addres, &u.OfficeUser.CreatedAt)
		if err != nil {
			return nil, err
		}
		newUsers = append(newUsers, u)
	}

	return newUsers, nil
}
