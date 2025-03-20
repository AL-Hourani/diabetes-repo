package main

import (
	"log"

	"github.com/AL-Hourani/care-center/api"
	"github.com/AL-Hourani/care-center/config"
	"github.com/AL-Hourani/care-center/data"
)


func main(){
	
	
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

    data.InitStorage(db)



	server := api.CreateNewAPIServer(config.GetEnv("PORT" , ""), db)
	if err := server.Run() ; err != nil {
		log.Fatal(err)
	}

	
}