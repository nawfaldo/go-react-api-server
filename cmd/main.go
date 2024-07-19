package main

import (
	"log"
	"test/cmd/api"
	"test/config"
	"test/db"

	"github.com/go-sql-driver/mysql"
)

func main() {
	_db, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}

	db.InitStorage(_db)

	server := api.NewAPIServer(":8080", _db)
	if err := server.Run(); err != nil {
		log.Fatal()
	}
}
