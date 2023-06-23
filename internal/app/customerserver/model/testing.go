package model

import (
	"testing"
)

func TestOffice(t *testing.T) *Office {
	return &Office{
		Name:   "OfficeName",
		Addres: "OfficeAddr",
	}
}

func TestUser(t *testing.T, of *Office) *User {
	return &User{
		Name:       "OfficeName",
		OfficeUser: of,
	}
}
