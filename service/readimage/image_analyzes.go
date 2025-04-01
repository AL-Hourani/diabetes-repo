package readimage

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"github.com/AL-Hourani/care-center/utils"
	"github.com/gorilla/mux"
)

type Handler struct {

}

func NewHandler( ) *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/upload",h.uploadHandler).Methods("POST")

}


func (h *Handler) uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to get image file", http.StatusBadRequest)
		return
	}
	defer file.Close()


	dst, err := os.Create(fmt.Sprintf("./uploads/%s", header.Filename))
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}
	defer dst.Close()


	if _, err = io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

    utils.WriteJSON(w , http.StatusOK ,[]byte("Image uploaded successfully"))
}

