package supervisor

import (
	"net/http"

	"github.com/AL-Hourani/care-center/service/auth"
	"github.com/AL-Hourani/care-center/types"
	"github.com/AL-Hourani/care-center/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.CenterStore
    pStore types.PatientStore
    superStore types.SuperisorStore
}

func NewHandler(store types.CenterStore , patientStore types.PatientStore , supervisorStore types.SuperisorStore) *Handler {
	return &Handler {
		store: store,
		pStore: patientStore,
	    superStore: supervisorStore,
	}
}




func (h *Handler) RegisterSuperVisorRoutes(router *mux.Router) {
	router.HandleFunc("/getSupervisorCenters",auth.WithJWTAuth(h.handleGetAllCentersData)).Methods("GET")
}







func (h *Handler) handleGetAllCentersData(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value(auth.UserContextKey).(*jwt.Token)
	if !ok {
		http.Error(w, "Unauthorized: No token found", http.StatusUnauthorized)
		return
	}


	id, err := auth.GetIDFromToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	user , err := h.pStore.GetLoginByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	if user.Role != "supervisor" {
		http.Error(w, "Unauthorized: You are not supervisor", http.StatusUnauthorized)
		return
	}

	centerData , err := h.superStore.GetAllCenters()
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err)
		return
	}

	utils.WriteJSON(w , http.StatusOK , centerData	)

}