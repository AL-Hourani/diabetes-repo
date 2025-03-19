package dashboard

import (
	"fmt"
	"net/http"

	"github.com/AL-Hourani/care-center/types"
	"github.com/AL-Hourani/care-center/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {

}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterDashboardRoutes(router *mux.Router) {
	router.HandleFunc("/RegisterPatientHealthOverview", h.handleHealthOverviewRegisteraion).Methods("POST")
}


func (h *Handler) handleHealthOverviewRegisteraion(w http.ResponseWriter , r *http.Request) {
	//get json payload
	var patientHealthOverview types.RegisterPatientHealthOverviewPayload
	if err := utils.ParseJSON(r , &patientHealthOverview); err != nil {
		utils.WriteError(w , http.StatusBadRequest , err)
	}
 
	// validate you info
	if err := utils.Validate.Struct(patientHealthOverview);err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w , http.StatusBadRequest , fmt.Errorf("invalid payload %v", error) )
		return
	}

	





}