package data

import (
	"database/sql"
	"log"
    "fmt"
	
	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}


func (cfg *PostgresConfig) FormatDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
}

func NewMQLStorage(cfg PostgresConfig) (*sql.DB , error) {
	db , err := sql.Open("postgres", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}


	db.SetMaxOpenConns(50)    // الحد الأقصى للاتصالات المفتوحة في نفس الوقت
    db.SetMaxIdleConns(20)     // الحد الأقصى للاتصالات الفارغة التي يمكن أن تبقى مفتوحة
    db.SetConnMaxLifetime(0)  

	return db , nil

}


func InitStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB : Successfuly connected !")
}