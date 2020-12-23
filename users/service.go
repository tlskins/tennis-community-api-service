package users

import (
	"log"

	m "github.com/tennis-community-api-service/pkg/mongo"
	"github.com/tennis-community-api-service/users/store"
)

type UsersService struct {
	Store *store.UsersStore
}

func Init(mongoDBName, mongoHost, mongoUser, mongoPwd string) (*UsersService, error) {
	mc, err := m.NewClientV2(mongoHost, mongoUser, mongoPwd)
	if err != nil {
		log.Fatalln(err)
	}
	newStore := store.NewStore(mc, mongoDBName)

	return &UsersService{Store: newStore}, nil
}
