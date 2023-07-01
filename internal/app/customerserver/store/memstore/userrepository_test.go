package memstore_test

import (
	"strconv"
	"testing"

	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/model"
	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/store/memstore"
	"github.com/stretchr/testify/assert"
)

func TestUserRep_Add(t *testing.T) {
	st := memstore.New()
	of := model.TestOffice(t)
	st.Office().Add(of)
	u := model.TestUser(t, of)
	err := st.User().Add(u)
	assert.NotNil(t, u)
	assert.NoError(t, err)
}

func TestUserRep_GetList(t *testing.T) {
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
			st.Office().Add(of)
			for i := 1; i <= tc.CountRow; i++ {
				u := model.TestUser(t, of)
				st.User().Add(u)
			}
			resp, _ := st.User().GetList()
			assert.Equal(t, tc.CountRow, len(resp))
		})
	}
}
