package service

import (
	"context"
	"log"

	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/restaurant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateRestNewMenuServiceClient(grpcClAddr string) restaurant.MenuServiceClient {
	conn, err := grpc.DialContext(context.Background(), grpcClAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to get log lvl, %v", err)
	}
	return restaurant.NewMenuServiceClient(conn)
}
