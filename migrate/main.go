package main

import (
	"log"
	"os"

	"github.com/AL-Hourani/care-center/config"
	"github.com/AL-Hourani/care-center/data"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	
	_ "github.com/golang-migrate/migrate/v4/source/file"
)



func main() {
	db , err := data.NewMQLStorage(data.PostgresConfig{
		User: 			  config.Envs.DBUser,
		Password:         config.Envs.DBPassword,
		Port:             config.Envs.DBPort,
		Host:    		  config.Envs.DBHost,
		DBName:           config.Envs.DBName,
		SSLMode: "require",
	})

	if err != nil {
		log.Fatal(err)
	}

    driver , err :=postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		log.Fatal(err)
	}

	m , err := migrate.NewWithDatabaseInstance(
		"file://migrate/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[(len(os.Args) - 1)]
	if cmd == "up" {
		if err := m.Up();err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	if cmd == "down" {
		if err := m.Down();err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

}