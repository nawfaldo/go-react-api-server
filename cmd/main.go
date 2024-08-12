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
		Passwd:               config.Envs.DBPwd,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		Addr:                 "sng101.hawkhost.com",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}

	db.InitStorage(_db)

	server := api.NewAPIServer(":8090", _db)
	if err := server.Run(); err != nil {
		log.Fatal()
	}
}
