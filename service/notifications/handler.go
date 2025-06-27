package notifications

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AL-Hourani/care-center/types"
)



type Hub struct {
	clients    map[int][]chan types.Notification
	register   chan clientConn
	unregister chan clientConn
	Broadcast  chan types.Notification
}

type clientConn struct {
	ID   int
	Chan chan types.Notification
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int][]chan types.Notification),
		register:   make(chan clientConn),
		unregister: make(chan clientConn),
		Broadcast:  make(chan types.Notification),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.ID] = append(h.clients[client.ID], client.Chan)

		case client := <-h.unregister:
			if chans, ok := h.clients[client.ID]; ok {
				for i, c := range chans {
					if c == client.Chan {
						h.clients[client.ID] = append(chans[:i], chans[i+1:]...)
						close(c)
						break
					}
				}
			}

		case notif := <-h.Broadcast:
			if chans, ok := h.clients[notif.ReceiverID]; ok {
				for _, ch := range chans {
					select {
					case ch <- notif:
					case <-time.After(1 * time.Second): // لا تنتظر كثيرًا
					}
				}
			}
		}
	}
}

func (h *Hub) HandleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Server does not support streaming", http.StatusInternalServerError)
		return
	}

	idStr := r.URL.Query().Get("patient_id")
	id, _ := strconv.Atoi(idStr)


	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// قناة لهذا العميل
	notifChan := make(chan types.Notification)

	// سجل العميل
	h.register <- clientConn{ID: id, Chan: notifChan}
	defer func() {
		h.unregister <- clientConn{ID: id, Chan: notifChan}
	}()

	// بث الإشعارات
	for notif := range notifChan {
		data, _ := json.Marshal(notif)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}
}
