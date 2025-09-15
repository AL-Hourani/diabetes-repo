package data

import (
	"database/sql"
	"log"
    "fmt"
	
	_ "github.com/lib/pq"
	//  "github.com/jackc/pgx/v5/pgxpool"
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


	db.SetMaxOpenConns(50)   
    db.SetMaxIdleConns(20)  
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





// type dbLogger struct{}


// func (d dbLogger) BeforeQuery(ctx context.Context, evt pgxpool.QueryEvent) context.Context {
//     f, _ := os.OpenFile("queries.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// fmt.Fprintln(f, "Executing query:", evt.SQL)
// f.Close()

// }


// func (d dbLogger) AfterQuery(ctx context.Context, evt pgxpool.QueryEvent) {
//     f, _ := os.OpenFile("queries.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// fmt.Fprintln(f, "Executing query:", evt.SQL)
// f.Close()

// }