package albums

import (
	"log"

	"github.com/tennis-community-api-service/moderation/store"
	m "github.com/tennis-community-api-service/pkg/mongo"
)

type ModerationService struct {
	Store *store.ModerationStore
}

func Init(mongoDBName, mongoHost, mongoUser, mongoPwd string) (*ModerationService, error) {
	mc, err := m.NewClientV2(mongoHost, mongoUser, mongoPwd)
	if err != nil {
		log.Fatalln(err)
	}
	newStore := store.NewStore(mc, mongoDBName)

	return &ModerationService{
		Store: newStore,
	}, nil
}
