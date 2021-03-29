package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"

	alb "github.com/tennis-community-api-service/albums"
	aT "github.com/tennis-community-api-service/albums/types"
	"github.com/tennis-community-api-service/pkg/enums"
)

func main() {
	cfgPath := flag.String("config", "../config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	albumsDBName := os.Getenv("ALBUMS_DB_NAME")
	albumsDBHost := os.Getenv("ALBUMS_DB_HOST")
	albumsDBUser := os.Getenv("ALBUMS_DB_USER")
	albumsDBPwd := os.Getenv("ALBUMS_DB_PWD")

	var albSvc *alb.AlbumsService
	var err error
	albSvc, err = alb.Init(albumsDBName, albumsDBHost, albumsDBUser, albumsDBPwd)
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	_, err = albSvc.CreateAlbum(context.TODO(), &aT.Album{
		Name:      "Federer Backhand",
		CreatedAt: now,
		UpdatedAt: now,
		Status:    enums.AlbumStatusCreated,
		IsPro:     true,
		SwingVideos: []*aT.SwingVideo{
			{
				ID:       uuid.NewV4().String(),
				Name:     "1",
				Frames:   210,
				VideoURL: "https://tennis-swings.s3.amazonaws.com/public/federer_backhand.mp4",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = albSvc.CreateAlbum(context.TODO(), &aT.Album{
		Name:      "Federer Forehand",
		CreatedAt: now,
		UpdatedAt: now,
		Status:    enums.AlbumStatusCreated,
		IsPro:     true,
		SwingVideos: []*aT.SwingVideo{
			{
				ID:       uuid.NewV4().String(),
				Name:     "1",
				Frames:   180,
				VideoURL: "https://tennis-swings.s3.amazonaws.com/public/federer_forehand.mp4",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = albSvc.CreateAlbum(context.TODO(), &aT.Album{
		Name:      "Federer Serve",
		CreatedAt: now,
		UpdatedAt: now,
		Status:    enums.AlbumStatusCreated,
		IsPro:     true,
		SwingVideos: []*aT.SwingVideo{
			{
				ID:       uuid.NewV4().String(),
				Name:     "1",
				Frames:   330,
				VideoURL: "https://tennis-swings.s3.amazonaws.com/public/federer_serve.mp4",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = albSvc.CreateAlbum(context.TODO(), &aT.Album{
		Name:      "Djokovic Backhand",
		CreatedAt: now,
		UpdatedAt: now,
		Status:    enums.AlbumStatusCreated,
		IsPro:     true,
		SwingVideos: []*aT.SwingVideo{
			{
				ID:       uuid.NewV4().String(),
				Name:     "1",
				Frames:   240,
				VideoURL: "https://tennis-swings.s3.amazonaws.com/public/djokovic_backhand.mp4",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = albSvc.CreateAlbum(context.TODO(), &aT.Album{
		Name:      "Djokovic Forehand",
		CreatedAt: now,
		UpdatedAt: now,
		Status:    enums.AlbumStatusCreated,
		IsPro:     true,
		SwingVideos: []*aT.SwingVideo{
			{
				ID:       uuid.NewV4().String(),
				Name:     "1",
				Frames:   330,
				VideoURL: "https://tennis-swings.s3.amazonaws.com/public/djokovic_forehand.mp4",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = albSvc.CreateAlbum(context.TODO(), &aT.Album{
		Name:      "Djokovic Serve",
		CreatedAt: now,
		UpdatedAt: now,
		Status:    enums.AlbumStatusCreated,
		IsPro:     true,
		SwingVideos: []*aT.SwingVideo{
			{
				ID:       uuid.NewV4().String(),
				Name:     "1",
				Frames:   540,
				VideoURL: "https://tennis-swings.s3.amazonaws.com/public/djokovic_serve.mp4",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = albSvc.CreateAlbum(context.TODO(), &aT.Album{
		Name:      "Nadal Backhand",
		CreatedAt: now,
		UpdatedAt: now,
		Status:    enums.AlbumStatusCreated,
		IsPro:     true,
		SwingVideos: []*aT.SwingVideo{
			{
				ID:       uuid.NewV4().String(),
				Name:     "1",
				Frames:   150,
				VideoURL: "https://tennis-swings.s3.amazonaws.com/public/nadal_backhand.mp4",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = albSvc.CreateAlbum(context.TODO(), &aT.Album{
		Name:      "Nadal Forehand",
		CreatedAt: now,
		UpdatedAt: now,
		Status:    enums.AlbumStatusCreated,
		IsPro:     true,
		SwingVideos: []*aT.SwingVideo{
			{
				ID:       uuid.NewV4().String(),
				Name:     "1",
				Frames:   450,
				VideoURL: "https://tennis-swings.s3.amazonaws.com/public/nadal_forehand.mp4",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = albSvc.CreateAlbum(context.TODO(), &aT.Album{
		Name:      "Nadal Serve",
		CreatedAt: now,
		UpdatedAt: now,
		Status:    enums.AlbumStatusCreated,
		IsPro:     true,
		SwingVideos: []*aT.SwingVideo{
			{
				ID:       uuid.NewV4().String(),
				Name:     "1",
				Frames:   450,
				VideoURL: "https://tennis-swings.s3.amazonaws.com/public/nadal_serve.mp4",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done")
}
