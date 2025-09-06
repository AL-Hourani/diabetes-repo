package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/AL-Hourani/care-center/service/center"
	"github.com/AL-Hourani/care-center/service/notifications"
	"github.com/AL-Hourani/care-center/service/patient"
	"github.com/AL-Hourani/care-center/service/readimage"
	"github.com/AL-Hourani/care-center/service/session"
	"github.com/AL-Hourani/care-center/service/supervisor"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	var notifHub = notifications.NewHub()
	go notifHub.Run()


	sessionManager := session.NewManager([]byte("fgfggfgfDDggjg$#jjrjr8733DDdffkfjf6363hhhdhddhddhd"))
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1/").Subrouter()

	
	// center .....
	centerStore := center.NewStore(s.db)
	patientStore := patient.NewStore(s.db)
	supervisorStore := supervisor.NewStore(s.db)

	centerHandler := center.NewHandler(centerStore ,patientStore , *sessionManager , notifHub , supervisorStore)
	centerHandler.RegisterCenterRoutes(subrouter)
	// patients ....

	//supervisor
	superVisorHandler := supervisor.NewHandler(centerStore ,patientStore , supervisorStore)
	superVisorHandler.RegisterSuperVisorRoutes(subrouter)

	patientHandler := patient.NewHandler(patientStore , centerStore ,*sessionManager , notifHub )
	patientHandler.RegisterPatientRoutes(subrouter)

	imageHandler := readimage.NewHandler()
	imageHandler.RegisterRoutes(subrouter)

	

   
		// إعدادات CORS
		cors := handlers.CORS(
			handlers.AllowedOrigins([]string{"https://mustafa-ki-o.github.io"}),
			handlers.AllowedMethods([]string{"GET", "POST","PATCH", "PUT", "DELETE", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		)


	
	log.Println("Listing on " , s.addr)
	return http.ListenAndServe(s.addr , cors(router))
}

