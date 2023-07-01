package model

import "github.com/google/uuid"

type Orders struct {
	UserUuid  uuid.UUID
	Salads    []OrderItem
	Garnishes []OrderItem
	Meats     []OrderItem
	Soups     []OrderItem
	Drinks    []OrderItem
	Desserts  []OrderItem
}
