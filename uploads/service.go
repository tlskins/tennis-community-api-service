package uploads

import (
	"log"

	m "github.com/tennis-community-api-service/pkg/mongo"
	"github.com/tennis-community-api-service/uploads/store"
)

type UploadsService struct {
	Store *store.UploadsStore
}

func Init(mongoDBName, mongoHost, mongoUser, mongoPwd string) (*UploadsService, error) {
	mc, err := m.NewClientV2(mongoHost, mongoUser, mongoPwd)
	if err != nil {
		log.Fatalln(err)
	}
	newStore := store.NewStore(mc, mongoDBName)

	return &UploadsService{Store: newStore}, nil
}
