package memstore

import (
	"sync"
	"time"

	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/model"
	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/store"
	"github.com/google/uuid"
)

type OfficeRepository struct {
	mu      sync.Mutex
	store   *Store
	offices map[uuid.UUID]*model.Office
}

func (r *OfficeRepository) Add(o *model.Office) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var err error
	o.CreatedAt = time.Now()
	o.Uuid, err = uuid.NewUUID()
	if err != nil {
		return err
	}
	r.offices[o.Uuid] = o
	return nil
}

func (r *OfficeRepository) GetList() []*model.Office {
	r.mu.Lock()
	defer r.mu.Unlock()
	newOffices := make([]*model.Office, 0, len(r.offices))
	for _, val := range r.offices {
		newOffices = append(newOffices, &model.Office{
			Uuid:      val.Uuid,
			Name:      val.Name,
			Addres:    val.Addres,
			CreatedAt: val.CreatedAt,
		})
	}
	return newOffices
}

func (r *OfficeRepository) Find(id uuid.UUID) (*model.Office, error) {
	office, ok := r.offices[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return office, nil
}
