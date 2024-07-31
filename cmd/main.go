package main

import (
	"log"
	"os"
	"test/cmd/api"
	"test/db"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	_db, err := db.NewPostgresStorage("postgres://default:Ut4uNix0wdRk@ep-polished-sea-a1efivnq.ap-southeast-1.aws.neon.tech:5432/verceldb?sslmode=require")
	if err != nil {
		log.Fatal(err)
	}

	db.InitStorage(_db)

	server := api.NewAPIServer(":4000", _db)
	if err := server.Run(); err != nil {
		log.Fatal()
	}
}
