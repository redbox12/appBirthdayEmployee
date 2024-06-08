package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"os"
	"taskRumbler/pkg/api"
	"taskRumbler/pkg/repository"
)

func main() {

	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	urlDb := os.Getenv("DATABASE_URL")
	//port := os.Getenv("PORT")
	//
	//if port == "" {
	//	port = "8000" //localhost
	//}
	//addr := ":" + port

	db, err := repository.NewPGRepo(urlDb)
	if err != nil {
		log.Fatal(err)
	}

	api := api.New(mux.NewRouter(), db)
	api.Handle()
	log.Fatal(api.ListenAndServe("localhost:8090"))

}
