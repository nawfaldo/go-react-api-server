package main

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// func main() {
// 	db, err := db.NewMySQLStorage(mysqlCfg.Config{
// 		User:                 config.Envs.DBUser,
// 		DBName:               config.Envs.DBName,
// 		Net:                  "tcp",
// 		AllowNativePasswords: true,
// 		ParseTime:            true,
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	driver, err := mysql.WithInstance(db, &mysql.Config{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	m, err := migrate.NewWithDatabaseInstance(
// 		"file://cmd/migrate/migrations",
// 		"mysql",
// 		driver,
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	cmd := os.Args[(len(os.Args) - 1)]
// 	if cmd == "up" {
// 		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
// 			log.Fatal(err)
// 		}
// 	}
// 	if cmd == "down" {
// 		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
// 			log.Fatal(err)
// 		}
// 	}
// }
