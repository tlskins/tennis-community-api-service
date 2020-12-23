package albums

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/tennis-community-api-service/albums/store"
	t "github.com/tennis-community-api-service/albums/types"
	m "github.com/tennis-community-api-service/pkg/mongo"

	"github.com/joho/godotenv"
)

type AlbumService interface {
	CreateSwingUpload(ctx context.Context, data *t.SwingUpload) (*t.SwingUpload, error)
	CreateUploadClipVideos(ctx context.Context, bucket string, outputs []string) (*t.SwingUpload, error)
	CreateUploadSwingVideos(ctx context.Context, bucket string, outputs []string) (*t.SwingUpload, error)
}

func Init() (AlbumService, error) {
	cfgPath := flag.String("config", "config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	mongoDBName := os.Getenv("MONGO_DB_NAME")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoUser := os.Getenv("MONGO_USER")
	mongoPwd := os.Getenv("MONGO_PWD")
	jwtKeyPath := os.Getenv("JWT_KEY_PATH")
	jwtSecretPath := os.Getenv("JWT_SECRET_PATH")

	mc, err := m.NewClientV2(mongoHost, mongoUser, mongoPwd)
	if err != nil {
		log.Fatalln(err)
	}
	storeInstance := store.NewStore(mc, mongoDBName)

	// jwtKey, err := ioutil.ReadFile(jwtKeyPath)
	// if err != nil {
	// 	return nil, err
	// }
	// jwtSecret, err := ioutil.ReadFile(jwtSecretPath)
	// if err != nil {
	// 	return nil, err
	// }
	// j, err := auth.NewJWTService(auth.JWTServiceConfig{
	// 	Key:    jwtKey,
	// 	Secret: jwtSecret,
	// 	RPCPwd: rpcPwd,
	// })
	// if err != nil {
	// 	nil, err
	// }

	return &Usecase{
		Store: storeInstance,
	}, nil
}
