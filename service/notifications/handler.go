package notifications

import (
    "net/http"
    "strconv"

    "github.com/gorilla/websocket"
)

type Notification struct {
    SenderID   int    `json:"sender_id"`
    ReceiverID int    `json:"receiver_id"`
    Message    string `json:"message"`
}

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // يسمح بالاتصال من أي origin (مطلوب لـ Frontend خارجي)
    },
}

type Hub struct {
    clients    map[int]*websocket.Conn
    register   chan clientConn
    unregister chan int
    Broadcast  chan Notification
}

type clientConn struct {
    ID   int
    Conn *websocket.Conn
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[int]*websocket.Conn),
        register:   make(chan clientConn),
        unregister: make(chan int),
        Broadcast:  make(chan Notification),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client.ID] = client.Conn
            // fmt.Println("✅ متصل جديد:", client.ID)

        case id := <-h.unregister:
            if conn, ok := h.clients[id]; ok {
                conn.Close()
                delete(h.clients, id)
            }

        case notif := <-h.Broadcast:
            if conn, ok := h.clients[notif.ReceiverID]; ok {
                conn.WriteJSON(notif)
            }
        }
    }
}

func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("patient_id")
    id, _ := strconv.Atoi(idStr)

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, "فشل الاتصال", 400)
        return
    }

    h.register <- clientConn{ID: id, Conn: conn}
    defer func() {
        h.unregister <- id
    }()

    // ننتظر الرسائل (حتى لو لم نستخدمها حاليًا)
    for {
        var dummy string
        if err := conn.ReadJSON(&dummy); err != nil {
            break
        }
    }
}
