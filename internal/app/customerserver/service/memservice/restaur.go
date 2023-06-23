package memservice

import (
	"context"

	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/restaurant"
	"google.golang.org/grpc"
)

type memRestMenuCL struct {
	restaurant.MenuServiceClient
	menu *restaurant.GetMenuResponse
}

func (m *memRestMenuCL) GetMenu(ctx context.Context, in *restaurant.GetMenuRequest, opts ...grpc.CallOption) (*restaurant.GetMenuResponse, error) {
	return m.menu, nil
}

func Create() *memRestMenuCL {
	m := &memRestMenuCL{}
	return m
}

func (m *memRestMenuCL) SetMenu(mn *restaurant.GetMenuResponse) {
	m.menu = mn
}

func InitMenuCl(m *memRestMenuCL) restaurant.MenuServiceClient {
	return m
}
