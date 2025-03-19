package main

import (
	"log"

	"github.com/AL-Hourani/care-center/api"
	"github.com/AL-Hourani/care-center/config"
	"github.com/AL-Hourani/care-center/data"
	"github.com/go-sql-driver/mysql"
)


func main(){
	
	
	db , err := data.NewMQLStorage(mysql.Config{
		User: 			config.Envs.DBUser,
		Passwd:         config.Envs.DBPassword,
		Addr:           config.Envs.DBAddress,
		DBName:         config.Envs.DBName,
		Net: "tcp",
		AllowNativePasswords: true,
		ParseTime: true,
	})

	if err != nil {
		log.Fatal(err)
	}

    data.InitStorage(db)



	server := api.CreateNewAPIServer(":8080", db)
	if err := server.Run() ; err != nil {
		log.Fatal(err)
	}

	
}