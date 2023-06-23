package customerserver

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/broker"
	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/model"
	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/store"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/customer"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/restaurant"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	customer.UnimplementedOfficeServiceServer
	customer.UnimplementedOrderServiceServer
	customer.UnimplementedUserServiceServer
	mux        *runtime.ServeMux
	logger     *logrus.Logger
	store      store.Store
	producer   sarama.SyncProducer
	servMenuCl restaurant.MenuServiceClient
}

func newServer(store store.Store, producer sarama.SyncProducer, servMenuCl restaurant.MenuServiceClient) *server {
	s := &server{
		mux:        runtime.NewServeMux(),
		logger:     logrus.New(),
		store:      store,
		producer:   producer,
		servMenuCl: servMenuCl,
	}

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.logger.Debugf("started %s %s", r.Method, r.RequestURI)
	start := time.Now()
	s.mux.ServeHTTP(w, r)
	s.logger.Debugf(
		"completed in %v",
		time.Since(start),
	)
}

func (s *server) CreateOffice(ctx context.Context, req *customer.CreateOfficeRequest) (*customer.CreateOfficeResponse, error) {
	s.logger.Debugln("CreateOffice")
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	office := &model.Office{
		Name:   req.Name,
		Addres: req.Address,
	}
	if err := s.store.Office().Add(office); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &customer.CreateOfficeResponse{}, nil
}

func (s *server) GetOfficeList(ctx context.Context, req *customer.GetOfficeListRequest) (*customer.GetOfficeListResponse, error) {
	s.logger.Debugln("GetOfficeList")
	offices := s.store.Office().GetList()
	retList := make([]*customer.Office, 0, len(offices))
	for _, val := range offices {
		retList = append(retList, &customer.Office{
			Uuid:      val.Uuid.String(),
			Name:      val.Name,
			Address:   val.Addres,
			CreatedAt: timestamppb.New(val.CreatedAt),
		})
	}

	return &customer.GetOfficeListResponse{Result: retList}, nil
}

func (s *server) CreateUser(ctx context.Context, req *customer.CreateUserRequest) (*customer.CreateUserResponse, error) {
	s.logger.Debugln("CreateUser")
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	office, err := s.store.Office().Find(uuid.MustParse(req.OfficeUuid))
	if err != nil {
		if err == store.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, err.Error())
		} else {
			return nil, status.Errorf(codes.Unknown, err.Error())
		}
	}

	user := &model.User{
		Name:       req.Name,
		OfficeUser: office,
	}

	if err := s.store.User().Add(user); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &customer.CreateUserResponse{}, nil
}

func (s *server) GetUserList(ctx context.Context, req *customer.GetUserListRequest) (*customer.GetUserListResponse, error) {
	s.logger.Debugln("GetUserList")
	users := s.store.User().GetList()
	retList := make([]*customer.User, 0, len(users))
	for _, val := range users {
		retList = append(retList, &customer.User{
			Uuid:       val.Uuid.String(),
			Name:       val.Name,
			OfficeUuid: val.OfficeUser.Uuid.String(),
			OfficeName: val.OfficeUser.Name,
			CreatedAt:  timestamppb.New(val.CreatedAt),
		})
	}

	return &customer.GetUserListResponse{Result: retList}, nil
}

func (s *server) CreateOrder(ctx context.Context, req *customer.CreateOrderRequest) (*customer.CreateOrderResponse, error) {
	s.logger.Debugln("CreateOrder")
	if _, err := uuid.Parse(req.UserUuid); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	msgObj := model.Orders{
		UserUuid: uuid.MustParse(req.UserUuid),
	}

	makeOrders := func(inOrders []*customer.OrderItem) []model.OrderItem {
		retOrdres := make([]model.OrderItem, 0, len(inOrders))
		for _, val := range inOrders {
			retOrdres = append(retOrdres, model.OrderItem{
				Count:       val.Count,
				ProductUuid: uuid.MustParse(val.ProductUuid),
			})
		}
		return retOrdres
	}

	msgObj.Salads = makeOrders(req.Salads)
	msgObj.Garnishes = makeOrders(req.Garnishes)
	msgObj.Meats = makeOrders(req.Meats)
	msgObj.Soups = makeOrders(req.Soups)
	msgObj.Drinks = makeOrders(req.Drinks)
	msgObj.Desserts = makeOrders(req.Desserts)

	msgStr, err := json.Marshal(msgObj)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	producerMsg := &sarama.ProducerMessage{Topic: broker.OrderCreatedTopic, Value: sarama.StringEncoder(msgStr)}
	_, _, err = s.producer.SendMessage(producerMsg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &customer.CreateOrderResponse{}, nil
}

func (s *server) GetActualMenu(ctx context.Context, req *customer.GetActualMenuRequest) (*customer.GetActualMenuResponse, error) {
	s.logger.Debugln("GetActualMenu")
	resp, err := s.servMenuCl.GetMenu(ctx, &restaurant.GetMenuRequest{})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	rspMenu := resp.Menu
	custmMenu := &customer.GetActualMenuResponse{}

	makeOrders := func(inOrders []*restaurant.Product) []*customer.Product {
		retOrdres := make([]*customer.Product, 0, len(inOrders))
		for _, val := range inOrders {
			retOrdres = append(retOrdres, &customer.Product{
				Uuid:        val.Uuid,
				Name:        val.Name,
				Description: val.Description,
				Type:        customer.CustomerProductType(val.Type),
				Weight:      val.Weight,
				Price:       val.Price,
				CreatedAt:   val.CreatedAt,
			})
		}
		return retOrdres
	}
	custmMenu.Salads = makeOrders(rspMenu.Salads)
	custmMenu.Garnishes = makeOrders(rspMenu.Garnishes)
	custmMenu.Meats = makeOrders(rspMenu.Meats)
	custmMenu.Soups = makeOrders(rspMenu.Soups)
	custmMenu.Drinks = makeOrders(rspMenu.Drinks)
	custmMenu.Drinks = makeOrders(rspMenu.Drinks)
	custmMenu.Desserts = makeOrders(rspMenu.Desserts)

	return custmMenu, nil
}
