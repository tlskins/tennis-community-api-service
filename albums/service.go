package albums

import (
	"log"

	"github.com/tennis-community-api-service/albums/store"
	m "github.com/tennis-community-api-service/pkg/mongo"
)

type AlbumsService struct {
	Store *store.AlbumsStore
}

func Init(mongoDBName, mongoHost, mongoUser, mongoPwd string) (*AlbumsService, error) {
	mc, err := m.NewClientV2(mongoHost, mongoUser, mongoPwd)
	if err != nil {
		log.Fatalln(err)
	}
	newStore := store.NewStore(mc, mongoDBName)

	return &AlbumsService{
		Store: newStore,
	}, nil
}
