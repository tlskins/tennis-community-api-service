package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"

	alb "github.com/tennis-community-api-service/albums"
	usr "github.com/tennis-community-api-service/users"
)

func main() {
	cfgPath := flag.String("config", "../config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	usersDBName := os.Getenv("USERS_DB_NAME")
	usersDBHost := os.Getenv("USERS_DB_HOST")
	usersDBUser := os.Getenv("USERS_DB_USER")
	usersDBPwd := os.Getenv("USERS_DB_PWD")
	albumsDBName := os.Getenv("ALBUMS_DB_NAME")
	albumsDBHost := os.Getenv("ALBUMS_DB_HOST")
	albumsDBUser := os.Getenv("ALBUMS_DB_USER")
	albumsDBPwd := os.Getenv("ALBUMS_DB_PWD")

	var usrSvc *usr.UsersService
	var albSvc *alb.AlbumsService
	var err error
	usrSvc, err = usr.Init(usersDBName, usersDBHost, usersDBUser, usersDBPwd)
	if err != nil {
		log.Fatal(err)
	}
	albSvc, err = alb.Init(albumsDBName, albumsDBHost, albumsDBUser, albumsDBPwd)
	if err != nil {
		log.Fatal(err)
	}

	if err = usrSvc.Store.EnsureIndexes(); err != nil {
		log.Fatal(err)
	}
	if err = albSvc.Store.EnsureIndexes(); err != nil {
		log.Fatal(err)
	}
	return
}
