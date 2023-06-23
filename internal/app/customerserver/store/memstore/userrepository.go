package memstore

import (
	"sync"
	"time"

	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/model"
	"github.com/google/uuid"
)

type UserRepository struct {
	mu    sync.Mutex
	store *Store
	users map[uuid.UUID]*model.User
}

func (r *UserRepository) Add(o *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var err error
	o.CreatedAt = time.Now()
	o.Uuid, err = uuid.NewUUID()
	if err != nil {
		return err
	}
	r.users[o.Uuid] = o
	return nil
}

func (r *UserRepository) GetList() []*model.User {
	r.mu.Lock()
	defer r.mu.Unlock()
	newUsers := make([]*model.User, 0, len(r.users))
	for _, val := range r.users {
		newUsers = append(newUsers, &model.User{
			Uuid: val.Uuid,
			Name: val.Name,
			OfficeUser: &model.Office{
				Uuid:      val.OfficeUser.Uuid,
				Name:      val.OfficeUser.Name,
				Addres:    val.OfficeUser.Addres,
				CreatedAt: val.OfficeUser.CreatedAt,
			},
			CreatedAt: val.CreatedAt,
		})
	}
	return newUsers
}
