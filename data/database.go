package data

import (
	"database/sql"
	"log"
	"github.com/go-sql-driver/mysql"
)



func NewMQLStorage(cfg mysql.Config) (*sql.DB , error) {
	db , err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	return db , nil
}

func InitStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB : Successfuly connected !")
}