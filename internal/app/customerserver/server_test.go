package customerserver

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/broker/membroker"
	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/model"
	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/service/memservice"
	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/store/memstore"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/customer"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/restaurant"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestServer_CreateOffice(t *testing.T) {
	s := newServer(memstore.New(), nil, nil)
	testCases := []struct {
		name    string
		office  *customer.CreateOfficeRequest
		isValid bool
	}{
		{
			name: "valid",
			office: &customer.CreateOfficeRequest{
				Name:    "NameTest",
				Address: "AddressTest",
			},
			isValid: true,
		},
		{
			name: "invalid name",
			office: &customer.CreateOfficeRequest{
				Name:    "Nam",
				Address: "AddressTest",
			},
			isValid: false,
		},
		{
			name: "invalid address",
			office: &customer.CreateOfficeRequest{
				Name:    "NameTest",
				Address: "Addr",
			},
			isValid: false,
		},
	}
	ctx := context.Background()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.CreateOffice(ctx, tc.office)
			if tc.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestServer_GetOfficeList(t *testing.T) {
	const ROW_COUNT int = 5
	s := newServer(memstore.New(), nil, nil)
	testCases := make([]struct {
		name     string
		CountRow int
	}, 0, ROW_COUNT)
	for i := 1; i < ROW_COUNT; i++ {
		testCases = append(testCases, struct {
			name     string
			CountRow int
		}{
			name:     "Count " + strconv.Itoa(i),
			CountRow: i,
		})
	}
	ctx := context.Background()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s.store.Office().Add(model.TestOffice(t))
			resp, _ := s.GetOfficeList(ctx, &customer.GetOfficeListRequest{})
			assert.Equal(t, tc.CountRow, len(resp.Result))
		})
	}
}

func TestServer_CreateUser(t *testing.T) {
	s := newServer(memstore.New(), nil, nil)
	office := model.TestOffice(t)
	s.store.Office().Add(office)
	testCases := []struct {
		name    string
		user    *customer.CreateUserRequest
		isValid bool
	}{
		{
			name: "valid",
			user: &customer.CreateUserRequest{
				Name:       "NameUser",
				OfficeUuid: office.Uuid.String(),
			},
			isValid: true,
		},
		{
			name: "invalid name",
			user: &customer.CreateUserRequest{
				Name:       "Nam",
				OfficeUuid: office.Uuid.String(),
			},
			isValid: false,
		},
		{
			name: "invalid OfficeUuid",
			user: &customer.CreateUserRequest{
				Name:       "NameTest",
				OfficeUuid: "ewr334",
			},
			isValid: false,
		},
		{
			name: "notFound Office",
			user: &customer.CreateUserRequest{
				Name:       "NameTest",
				OfficeUuid: uuid.NewString(),
			},
			isValid: false,
		},
	}
	ctx := context.Background()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.CreateUser(ctx, tc.user)
			if tc.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestServer_GetUserList(t *testing.T) {
	const ROW_COUNT int = 5
	s := newServer(memstore.New(), nil, nil)
	testCases := make([]struct {
		name     string
		CountRow int
	}, 0, ROW_COUNT)
	for i := 1; i < ROW_COUNT; i++ {
		testCases = append(testCases, struct {
			name     string
			CountRow int
		}{
			name:     "Count " + strconv.Itoa(i),
			CountRow: i,
		})
	}
	ctx := context.Background()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			office := model.TestOffice(t)
			s.store.Office().Add(office)
			s.store.User().Add(model.TestUser(t, office))
			resp, _ := s.GetUserList(ctx, &customer.GetUserListRequest{})
			assert.Equal(t, tc.CountRow, len(resp.Result))
		})
	}
}

func TestServer_CreateOrder(t *testing.T) {
	makeOrders := func(qnt int, prUuid string) []*customer.OrderItem {
		var prodUuid string
		if prUuid == "" {
			prodUuid = uuid.NewString()
		} else {
			prodUuid = prUuid
		}
		retOrdres := make([]*customer.OrderItem, 0, qnt)
		for i := 1; i <= qnt; i++ {
			retOrdres = append(retOrdres, &customer.OrderItem{
				Count:       int32(qnt),
				ProductUuid: prodUuid,
			})
		}
		return retOrdres
	}

	makeCreOrdReq := func() *customer.CreateOrderRequest {
		reqOrder := &customer.CreateOrderRequest{}
		reqOrder.UserUuid = uuid.NewString()
		reqOrder.Salads = makeOrders(2, "")
		reqOrder.Garnishes = makeOrders(1, "")
		reqOrder.Meats = makeOrders(3, "")
		reqOrder.Soups = makeOrders(2, "")
		reqOrder.Drinks = makeOrders(1, "")
		reqOrder.Desserts = makeOrders(5, "")
		return reqOrder
	}
	mp := membroker.Create()
	producer := membroker.IniProducer(mp)
	s := newServer(memstore.New(), producer, nil)
	testCases := []struct {
		name           string
		order          *customer.CreateOrderRequest
		checkOrderItem func(*testing.T, []*customer.OrderItem, []model.OrderItem)
		isValid        bool
	}{
		{
			name:  "valid",
			order: makeCreOrdReq(),
			checkOrderItem: func(t *testing.T, custOrdrIt []*customer.OrderItem, msgOrdIt []model.OrderItem) {
				for i := 0; i < len(msgOrdIt); i++ {
					assert.Equal(t, custOrdrIt[i].Count, msgOrdIt[i].Count)
					assert.Equal(t, custOrdrIt[i].ProductUuid, msgOrdIt[i].ProductUuid.String())
				}
			},
			isValid: true,
		},
		{
			name: "invalid UserUuid",
			order: &customer.CreateOrderRequest{
				UserUuid: "sdfd",
			},
			isValid: false,
		},
		/*{
			name: "invalid Salads ProductUuid",
			order: &customer.CreateOrderRequest{
				UserUuid: uuid.NewString(),
				Salads:   makeOrders(1, "esdd"),
			},
			isValid: false,
		},*/
	}
	ctx := context.Background()
	var (
		vlStr  []byte
		msgObj model.Orders
	)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.CreateOrder(ctx, tc.order)
			if tc.isValid {
				assert.NoError(t, err)
				vlStr, _ = mp.Msg.Value.Encode()
				_ = json.Unmarshal(vlStr, &msgObj)
				assert.Equal(t, tc.order.UserUuid, msgObj.UserUuid.String())
				assert.Equal(t, len(tc.order.Salads), len(msgObj.Salads))
				tc.checkOrderItem(t, tc.order.Salads, msgObj.Salads)
				assert.Equal(t, len(tc.order.Garnishes), len(msgObj.Garnishes))
				tc.checkOrderItem(t, tc.order.Garnishes, msgObj.Garnishes)
				assert.Equal(t, len(tc.order.Meats), len(msgObj.Meats))
				tc.checkOrderItem(t, tc.order.Meats, msgObj.Meats)
				assert.Equal(t, len(tc.order.Soups), len(msgObj.Soups))
				tc.checkOrderItem(t, tc.order.Soups, msgObj.Soups)
				assert.Equal(t, len(tc.order.Drinks), len(msgObj.Drinks))
				tc.checkOrderItem(t, tc.order.Drinks, msgObj.Drinks)
				assert.Equal(t, len(tc.order.Desserts), len(msgObj.Desserts))
				tc.checkOrderItem(t, tc.order.Desserts, msgObj.Desserts)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestServer_GetActualMenu(t *testing.T) {
	makeProds := func(name string, typePr restaurant.ProductType, qnt int) []*restaurant.Product {
		retProd := make([]*restaurant.Product, 0, 1)
		for i := 1; i <= qnt; i++ {
			retProd = append(retProd, &restaurant.Product{
				Uuid:        uuid.NewString(),
				Name:        name,
				Description: name,
				Type:        typePr,
				Weight:      int32(qnt * 100),
				Price:       float64(qnt + 100),
				CreatedAt:   timestamppb.New(time.Now()),
			})
		}
		return retProd
	}

	makeGetMenuResp := func() *restaurant.GetMenuResponse {
		retMenu := &restaurant.GetMenuResponse{Menu: &restaurant.Menu{}}
		retMenu.Menu.Salads = makeProds("Salads", restaurant.ProductType_PRODUCT_TYPE_SALAD, 2)
		retMenu.Menu.Garnishes = makeProds("Garnishes", restaurant.ProductType_PRODUCT_TYPE_GARNISH, 3)
		retMenu.Menu.Meats = makeProds("Meats", restaurant.ProductType_PRODUCT_TYPE_MEAT, 1)
		retMenu.Menu.Soups = makeProds("Soups", restaurant.ProductType_PRODUCT_TYPE_SOUP, 4)
		retMenu.Menu.Drinks = makeProds("Drinks", restaurant.ProductType_PRODUCT_TYPE_DRINK, 5)
		retMenu.Menu.Desserts = makeProds("Desserts", restaurant.ProductType_PRODUCT_TYPE_DESSERT, 2)
		return retMenu
	}
	ms := memservice.Create()
	restGetMenu := makeGetMenuResp()
	restMenu := restGetMenu.Menu
	ms.SetMenu(restGetMenu)
	menuSer := memservice.InitMenuCl(ms)
	s := newServer(memstore.New(), nil, menuSer)
	testCases := []struct {
		name      string
		checkProd func(*testing.T, []*restaurant.Product, []*customer.Product)
	}{
		{
			name: "valid",
			checkProd: func(t *testing.T, restProd []*restaurant.Product, custProd []*customer.Product) {
				for i := 0; i < len(restProd); i++ {
					assert.Equal(t, restProd[i].Uuid, custProd[i].Uuid)
					assert.Equal(t, restProd[i].Name, custProd[i].Name)
					assert.Equal(t, restProd[i].Description, custProd[i].Description)
					assert.Equal(t, restProd[i].Type.Number(), custProd[i].Type.Number())
					assert.Equal(t, restProd[i].Weight, custProd[i].Weight)
					assert.Equal(t, restProd[i].Price, custProd[i].Price)
					assert.Equal(t, restProd[i].CreatedAt, custProd[i].CreatedAt)
				}
			},
		},
	}
	ctx := context.Background()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			custMenu, err := s.GetActualMenu(ctx, &customer.GetActualMenuRequest{})
			assert.NoError(t, err)
			assert.Equal(t, len(restMenu.Salads), len(custMenu.Salads))
			tc.checkProd(t, restMenu.Salads, custMenu.Salads)
			assert.Equal(t, len(restMenu.Garnishes), len(custMenu.Garnishes))
			tc.checkProd(t, restMenu.Garnishes, custMenu.Garnishes)
			assert.Equal(t, len(restMenu.Meats), len(custMenu.Meats))
			tc.checkProd(t, restMenu.Meats, custMenu.Meats)
			assert.Equal(t, len(restMenu.Soups), len(custMenu.Soups))
			tc.checkProd(t, restMenu.Soups, custMenu.Soups)
			assert.Equal(t, len(restMenu.Drinks), len(custMenu.Drinks))
			tc.checkProd(t, restMenu.Drinks, custMenu.Drinks)
			assert.Equal(t, len(restMenu.Desserts), len(custMenu.Desserts))
			tc.checkProd(t, restMenu.Desserts, custMenu.Desserts)
		})
	}
}
