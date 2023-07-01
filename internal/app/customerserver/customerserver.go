package customerserver

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/broker"
	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/service"
	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver/store/pgstore"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/customer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Start(cfg *Config) error {
	s := grpc.NewServer()

	db, err := newDB(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	store := pgstore.New(db)
	producer := broker.IniProducer(cfg.BrokerAddr)
	srvMenu := service.CreateRestNewMenuServiceClient(cfg.GrpcClAddr)
	gs := newServer(store, producer, srvMenu)
	gs.mux = runtime.NewServeMux()

	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatalf("failed to get log lvl, %v", err)
	}
	gs.logger.SetLevel(level)
	ctx, cancel := context.WithCancel(context.Background())

	go runGRPCServer(cfg, s, gs)
	go runHTTPServer(ctx, cfg, gs)

	gracefulShutDown(s, cancel, db)

	return nil
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func runGRPCServer(cfg *Config, s *grpc.Server, gs *server) {
	customer.RegisterOfficeServiceServer(s, gs)
	customer.RegisterUserServiceServer(s, gs)
	customer.RegisterOrderServiceServer(s, gs)

	l, err := net.Listen("tcp", cfg.GrpcAddr)
	if err != nil {
		log.Fatalf("failed to listen tcp %s, %v", cfg.GrpcAddr, err)
	}

	log.Printf("starting listening grpc server at %s", cfg.GrpcAddr)
	if err := s.Serve(l); err != nil {
		log.Fatalf("error service grpc server %v", err)
	}
}

func runHTTPServer(ctx context.Context, cfg *Config, gs *server) {
	regServiceHand := func(rgFun func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error) {
		err := rgFun(
			ctx,
			gs.mux,
			"0.0.0.0"+cfg.GrpcAddr,
			[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		)
		if err != nil {
			log.Fatal(err)
		}
	}
	regServiceHand(customer.RegisterOfficeServiceHandlerFromEndpoint)
	regServiceHand(customer.RegisterUserServiceHandlerFromEndpoint)
	regServiceHand(customer.RegisterOrderServiceHandlerFromEndpoint)

	log.Printf("starting listening http server at %s", cfg.HttpAddr)
	if err := http.ListenAndServe(cfg.HttpAddr, gs); err != nil {
		log.Fatalf("error service http server %v", err)
	}
}

func gracefulShutDown(s *grpc.Server, cancel context.CancelFunc, db *sql.DB) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer db.Close()
	defer signal.Stop(ch)

	sig := <-ch
	errorMessage := fmt.Sprintf("%s %v - %s", "Received shutdown signal:", sig, "Graceful shutdown done")
	log.Println(errorMessage)
	s.GracefulStop()
	cancel()
}
