package memstore_test

import (
	"strconv"
	"testing"

	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/model"
	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/store/memstore"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestOfficeRep_Add(t *testing.T) {
	st := memstore.New()
	of := model.TestOffice(t)
	err := st.Office().Add(of)
	assert.NotNil(t, of)
	assert.NoError(t, err)
}

func TestOfficeRep_GetList(t *testing.T) {
	const ROW_COUNT int = 5
	of := model.TestOffice(t)
	testCases := make([]struct {
		name     string
		CountRow int
	}, 0, ROW_COUNT)
	for i := 1; i <= ROW_COUNT; i++ {
		testCases = append(testCases, struct {
			name     string
			CountRow int
		}{
			name:     "Count " + strconv.Itoa(i),
			CountRow: i,
		})
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			st := memstore.New()
			for i := 1; i <= tc.CountRow; i++ {
				st.Office().Add(of)
			}
			resp, _ := st.Office().GetList()
			assert.Equal(t, tc.CountRow, len(resp))
		})
	}
}

func TestOfficeRep_Find(t *testing.T) {
	st := memstore.New()
	of := model.TestOffice(t)
	st.Office().Add(of)

	testCases := []struct {
		name string
		find func(*testing.T, uuid.UUID)
		uid  uuid.UUID
	}{
		{
			name: "valid",
			find: func(t *testing.T, ud uuid.UUID) {
				retOf, err := st.Office().Find(ud)
				assert.NotNil(t, of)
				assert.NoError(t, err)
				assert.Equal(t, of.Uuid, retOf.Uuid)
			},
			uid: of.Uuid,
		},
		{
			name: "invalid",
			find: func(t *testing.T, ud uuid.UUID) {
				_, err := st.Office().Find(ud)
				assert.Error(t, err)
			},
			uid: uuid.New(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.find(t, tc.uid)
		})
	}
}
