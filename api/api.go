package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/AL-Hourani/care-center/service/center"
	"github.com/AL-Hourani/care-center/service/patient"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

type APIServer struct {
	addr	string
	db 		*sql.DB
}


func CreateNewAPIServer(addr string  , db *sql.DB) *APIServer {
	return &APIServer{
		addr:addr,
		db:db,
	}
}


func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1/").Subrouter()

	
	// center .....
	centerStore := center.NewStore(s.db)
	centerHandler := center.NewHandler(centerStore)
	centerHandler.RegisterCenterRoutes(subrouter)
	// patients ....

	patientStore := patient.NewStore(s.db)
	patientHandler := patient.NewHandler(patientStore , centerStore)
	patientHandler.RegisterPatientRoutes(subrouter)
   
		// إعدادات CORS
		cors := handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}), // السماح بجميع النطاقات
			handlers.AllowedMethods([]string{"GET", "POST","PATCH", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		)


	
	log.Println("Listing on " , s.addr)
	return http.ListenAndServe(s.addr , cors(router))
}

