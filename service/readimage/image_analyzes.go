package readimage

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"github.com/AL-Hourani/care-center/utils"
	"github.com/gorilla/mux"
	"regexp"
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

	sanitizedFilename := sanitizeFilename(header.Filename)

	// التأكد من وجود المجلد uploads
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(fmt.Sprintf("./uploads/%s",sanitizedFilename ))
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


func sanitizeFilename(filename string) string {
	// إزالة الأقواس والمسافات
	re := regexp.MustCompile(`[^a-zA-Z0-9\-_\.]`)
	return re.ReplaceAllString(filename, "_")
}