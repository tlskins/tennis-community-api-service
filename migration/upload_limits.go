package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	usr "github.com/tennis-community-api-service/users"
	uT "github.com/tennis-community-api-service/users/types"
)

func main() {
	cfgPath := flag.String("config", "../config.dev.yml", "path for yaml config")
	flag.Parse()
	godotenv.Load(*cfgPath)

	usersDBName := os.Getenv("USERS_DB_NAME")
	usersDBHost := os.Getenv("USERS_DB_HOST")
	usersDBUser := os.Getenv("USERS_DB_USER")
	usersDBPwd := os.Getenv("USERS_DB_PWD")

	var usrSvc *usr.UsersService
	var err error
	usrSvc, err = usr.Init(usersDBName, usersDBHost, usersDBUser, usersDBPwd)
	if err != nil {
		log.Fatal(err)
	}

	users, err := usrSvc.Store.GetAllUsers()
	if err != nil {
		log.Fatal(err)
	}

	limit := 5
	for _, user := range users {
		_, err = usrSvc.Store.UpdateUser(&uT.UpdateUser{
			ID:                 user.ID,
			WeeklyUploadsLimit: &limit,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Updated %s...\n", user.UserName)
	}
	return
}
