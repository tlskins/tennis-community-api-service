package uploads

import (
	"log"

	m "github.com/tennis-community-api-service/pkg/mongo"
	"github.com/tennis-community-api-service/uploads/store"
)

type UploadsService struct {
	Store     *store.UploadsStore
	awsConfig *AWSConfig
}

type AWSConfig struct {
	accessKeyId     string
	secretAccessKey string
	bucketName      string
}

func Init(mongoDBName, mongoHost, mongoUser, mongoPwd, awsAccessKeyId, awsSecretAccessKey, bucketName string) (*UploadsService, error) {
	mc, err := m.NewClientV2(mongoHost, mongoUser, mongoPwd)
	if err != nil {
		log.Fatalln(err)
	}
	newStore := store.NewStore(mc, mongoDBName)

	return &UploadsService{
		Store: newStore,
		awsConfig: &AWSConfig{
			accessKeyId:     awsAccessKeyId,
			secretAccessKey: awsSecretAccessKey,
			bucketName:      bucketName,
		},
	}, nil
}
